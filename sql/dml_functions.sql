CREATE OR REPLACE FUNCTION move_module(module_id UUID, new_parent_id UUID) 
RETURNS void AS $$
DECLARE 
old_parent UUID;
root_siblings UUID[];
old_root_id UUID;
new_root_id UUID;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  PERFORM increment_topic_version(module_id);
  select module_parents.parent_id into old_parent from module_parents where child_id = module_id;
  if old_parent is null then --this means, the moving module is currently the root module.
    select array_agg(child_id::UUID) into root_siblings from module_parents where parent_id = module_id;
    -- this is an edge case where the moving module is the root module and it has more than one direct descendant.
    if array_length(root_siblings,1) > 1 then
      raise notice 'first sibling is %', root_siblings[1];
      --there can only be one root module. this update makes all siblings of the new root module to its descendants.    
      update module_parents set parent_id = root_siblings[1] where parent_id = module_id and child_id <> root_siblings[1]; --array index starts at 1...
    end if;
    -- after that we promote the remaining module to the new root.
    delete from module_parents where parent_id = module_id;
    insert into module_parents values(module_id, new_parent_id);
  else
    update module_parents set parent_id = old_parent where parent_id = module_id;
    delete from module_parents where child_id = module_id;
    if new_parent_id is null then -- the moving module should be the new root.
      select module_trees.id into old_root_id from module_trees where topic_id = (select topic_id from modules where id = module_id) AND level = 0;
      insert into module_parents values(old_root_id, module_id); -- the old root is now its first child.
    else
      insert into module_parents values(module_id,new_parent_id);
    end if;
  end if;
  REFRESH MATERIALIZED VIEW module_trees;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION move_module_tree(module_id UUID, new_parent_id UUID) 
RETURNS void AS $$
DECLARE 
parent_level int;
module_level int;
old_root_id UUID;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  PERFORM increment_topic_version(module_id);
  if new_parent_id is NULL then
    select module_trees.id into old_root_id from module_trees where topic_id = (select topic_id from modules where id = module_id) AND level = 0;
    insert into module_parents values(old_root_id, module_id); -- the old root is now its first child.
    delete from module_parents where child_id = module_id;
  else
    if exists( select 1 from (select unnest (paths) as paths from module_trees where id = new_parent_id )t where paths like '%'||module_id||'%') then
      raise notice 'found path';
      select module_trees.level into module_level from module_trees where id = module_id;
      select module_trees.level into parent_level from module_trees where id = new_parent_id;
      if parent_level > module_level then
        RAISE EXCEPTION 'Parent-Child cyclus. moduleID --> % parentID -->%', module_id, new_parent_id
        USING HINT = 'Please check the levels of module and new parent.';
      end if;
    end if;
    delete from module_parents where child_id = module_id;
    insert into module_parents values(module_id,new_parent_id);
  end if;
  REFRESH MATERIALIZED VIEW module_trees;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION delete_module(module_id UUID) 
RETURNS void AS $$
DECLARE 
parent_level int;
module_level int;
old_root_id UUID;
old_parent UUID;
root_siblings UUID[];
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  PERFORM increment_topic_version(module_id);
  select module_parents.parent_id into old_parent from module_parents where child_id = module_id;
  if old_parent is null then 
    select array_agg(child_id::UUID) into root_siblings from module_parents where parent_id = module_id;
    if array_length(root_siblings,1) > 1 then
      update module_parents set parent_id = root_siblings[1] where parent_id = module_id and child_id <> root_siblings[1];
    end if;
  else
    update module_parents set parent_id = old_parent where parent_id = module_id;
  end if;
  delete from modules where id = module_id;
  REFRESH MATERIALIZED VIEW module_trees;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_module_tree(module_id UUID) 
RETURNS void AS $$
DECLARE 
topic_id UUID;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  PERFORM increment_topic_version(module_id);
  delete from modules where  id in (select id from (select id, unnest (paths) as paths from module_trees )t where paths like '%'||module_id||'%');
  REFRESH MATERIALIZED VIEW module_trees;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_module(id UUID, topic_id UUID, description text,  video_id UUID, script_id UUID,parent_ids VARIADIC UUID[]) 
RETURNS void AS $$
DECLARE
parent_id UUID;
BEGIN
  insert into modules (id,topic_id,description,video_id,script_id,version) values(id,topic_id,description,video_id,script_id,1);
  foreach parent_id in ARRAY parent_ids
  loop
    insert into module_parents(child_id, parent_id) values(id,parent_id);
  end loop;
  PERFORM increment_topic_version(id);
  REFRESH MATERIALIZED VIEW module_trees;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_topic_version(module_id UUID) 
RETURNS void AS $$
BEGIN
  update topics t
  set version =  t.version + 1 
  FROM modules m
  where m.topic_id = m.id AND m.id = module_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_version(topic_id UUID, assumed_version bigint) 
RETURNS void AS $$
DECLARE 
real_version bigint;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  SELECT t.version into real_version FROM topics t where t.id = topic_id;
  if real_version > assumed_version then
    RAISE EXCEPTION 'record was modified';
  end if;
END;
$$ LANGUAGE plpgsql;







