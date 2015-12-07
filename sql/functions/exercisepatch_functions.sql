\c lecture
--TODO unit test
drop function add_hint(UUID,UUID, UUID,int,text,int);
CREATE OR REPLACE FUNCTION add_hint(in_id UUID,in_task_id UUID, in_exercise_id UUID,in_position int, in_content text, in_cost int) 
RETURNS void AS $$
DECLARE
max_position int;
BEGIN
  PERFORM check_exercise_context(in_exercise_id,in_task_id);
  SET CONSTRAINTS ALL DEFERRED;
  select max(position) into max_position from hints where task_id = in_task_id;
  if max_position > in_position then
    update hints set position = position+1 where task_id = in_task_id AND position >= in_position;
    insert into hints(id,task_id,position,content,cost) values(in_id,in_task_id,in_position,in_content,in_cost);
  else 
    insert into hints(id,task_id,position,content,cost) values(in_id,in_task_id,max_position+1,in_content,in_cost);
  end if;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function remove_hint(UUID, UUID,int);
CREATE OR REPLACE FUNCTION remove_hint(in_task_id UUID, in_exercise_id UUID,in_position int) 
RETURNS void AS $$
BEGIN
  PERFORM check_exercise_context(in_exercise_id,in_task_id);
  SET CONSTRAINTS ALL DEFERRED;
  delete from hints where task_id = in_task_id AND position = in_position;
  update hints set position = position-1 where task_id = in_task_id AND position >= in_position;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function move_hint(UUID, UUID,int,int);
CREATE OR REPLACE FUNCTION move_hint(in_task_id UUID, in_exercise_id UUID,in_old_position int, in_new_position int) 
RETURNS void AS $$
DECLARE
var_id UUID;
BEGIN
  PERFORM check_exercise_context(in_exercise_id,in_task_id);
  SET CONSTRAINTS ALL DEFERRED;
  select id into var_id from hints where task_id = in_task_id AND position = in_old_position;
  if var_id is null then 
    RAISE EXCEPTION 'unknown hint position';
  end if;
  if in_new_position > in_old_position then
  update hints set position = position - 1 where task_id = in_task_id AND position  > in_old_position AND position <= in_new_position ; 
else
  update hints set position = position + 1 where task_id = in_task_id AND position < in_old_position AND position >= in_new_position ; 
  end if;
  update hints set position = in_new_position where id = var_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;


--TODO unit test
drop function replace_hint_content(UUID,UUID, int,text);
CREATE OR REPLACE FUNCTION replace_hint_content(in_task_id UUID, in_exercise_id UUID,in_position int, in_content text) 
RETURNS void AS $$
BEGIN
  PERFORM check_exercise_context(in_exercise_id,in_task_id);
  update hints set content = in_content where task_id = in_task_id AND position  = in_position;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function replace_hint_cost(UUID, UUID,int,int);
CREATE OR REPLACE FUNCTION replace_hint_cost(in_task_id UUID, in_exercise_id UUID, in_position int, in_cost int) 
RETURNS void AS $$
BEGIN
  PERFORM check_exercise_context(in_exercise_id,in_task_id);
  update hints set cost = in_cost where task_id = in_task_id AND position  = in_position;
END;
$$ LANGUAGE plpgsql;


--TODO unit test
drop function replace_task_content(UUID, int,text);
CREATE OR REPLACE FUNCTION replace_task_content(in_exercise_id UUID, in_position int, in_content text) 
RETURNS void AS $$
  update tasks set content = in_content where exercise_id = in_exercise_id AND position  = in_position;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;

--TODO unit test
drop function add_task(UUID,UUID, int,text);
CREATE OR REPLACE FUNCTION add_task(in_id UUID,in_exercise_id UUID, in_position int, in_content text) 
RETURNS void AS $$
DECLARE
max_position int;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  select max(position) into max_position from tasks where exercise_id = in_exercise_id;
  if max_position > in_position then
    update tasks set position = position+1 where exercise_id = in_exercise_id AND position >= in_position;
    insert into tasks(id,exercise_id,position,content) values(in_id,in_exercise_id,in_position,in_content);
  else 
    insert into tasks(id,exercise_id,position,content) values(in_id,in_exercise_id,max_position+1,in_content);
  end if;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function remove_task(UUID, int);
CREATE OR REPLACE FUNCTION remove_task(in_exercise_id UUID, in_position int) 
RETURNS void AS $$
  SET CONSTRAINTS ALL DEFERRED;
  delete from tasks where exercise_id = in_exercise_id AND position = in_position;
  update tasks set position = position-1 where exercise_id = in_exercise_id AND position >= in_position;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;

--TODO unit test
drop function move_task(UUID, int,int);
CREATE OR REPLACE FUNCTION move_task(in_exercise_id UUID, in_old_position int, in_new_position int) 
RETURNS void AS $$
 DECLARE
var_content text;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  select content into var_content from tasks where exercise_id = in_exercise_id AND position = in_old_position;
  if var_content is null then 
    RAISE EXCEPTION 'unknown task position';
  end if;
  if in_new_position > in_old_position then
  update tasks set position = position - 1 where exercise_id = in_exercise_id AND position  > in_old_position AND position <= in_new_position ; 
else
  update tasks set position = position + 1 where exercise_id = in_exercise_id AND position < in_old_position AND position >= in_new_position ; 
  end if;
  update tasks set position = in_new_position where position = in_old_position AND exercise_id = in_exercise_id AND content = var_content; 
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;


drop function check_exercise_context(UUID,UUID);
CREATE OR REPLACE FUNCTION check_exercise_context(in_context_id UUID, in_task_id UUID) 
RETURNS void AS $$
BEGIN
  if NOT exists(select 1 from tasks t where t.id = in_task_id AND t.exercise_id = in_context_id) then
    RAISE EXCEPTION 'Operation out of scope.';
  end if;
END;
$$ LANGUAGE plpgsql;
