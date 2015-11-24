drop materialized view if exists module_trees;

create materialized view module_trees as 
with recursive p(id,level,path,description,topic_id) as (
  select id,0, '/'||id ,description ,topic_id from modules m left join module_parents mp on m.id = mp.child_id where mp.child_id is null
  union
  select m.id,level+1, path || '/' || m.id, m.description, m.topic_id
  from modules m left join module_parents mp on m.id = mp.child_id inner join p on p.id = mp.parent_id
) select id, level, array_agg(path) as paths, description , topic_id from p group by p.id,p.level,p.description, p.topic_id order by level;

