\c lecture
--TODO unit test
drop function add_hint(UUID,UUID, int,text,int);
CREATE OR REPLACE FUNCTION add_hint(in_id UUID,in_exercise_id UUID, in_position int, in_content text, in_cost int) 
RETURNS void AS $$
DECLARE
max_position int;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  select max(position) into max_position from hints where exercise_id = in_exercise_id;
  if max_position > in_position then
    update hints set position = position+1 where exercise_id = in_exercise_id AND position >= in_position;
    insert into hints(id,exercise_id,position,content,cost) values(in_id,in_exercise_id,in_position,in_content,in_cost);
  else 
    insert into hints(id,exercise_id,position,content,cost) values(in_id,in_exercise_id,max_position+1,in_content,in_cost);
  end if;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function remove_hint(UUID, int);
CREATE OR REPLACE FUNCTION remove_hint(in_exercise_id UUID, in_position int) 
RETURNS void AS $$
  SET CONSTRAINTS ALL DEFERRED;
  delete from hints where exercise_id = in_exercise_id AND position = in_position;
  update hints set position = position-1 where exercise_id = in_exercise_id AND position >= in_position;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function move_hint(UUID, int,int);
CREATE OR REPLACE FUNCTION move_hint(in_exercise_id UUID, in_old_position int, in_new_position int) 
RETURNS void AS $$
  SET CONSTRAINTS ALL DEFERRED;
  update hints set position = in_new_position where exercise_id = in_exercise_id AND position = in_old_position;
  update hints set position = in_old_position where exercise_id = in_exercise_id AND position = in_new_position; 
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function replace_hint_content(UUID, int,text);
CREATE OR REPLACE FUNCTION replace_hint_content(in_exercise_id UUID, in_position int, in_content text) 
RETURNS void AS $$
  update hints set content = in_content where exercise_id = in_exercise_id AND position  = in_position;
$$ LANGUAGE sql;

--TODO unit test
drop function replace_hint_cost(UUID, int,int);
CREATE OR REPLACE FUNCTION replace_hint_cost(in_exercise_id UUID, in_position int, in_cost int) 
RETURNS void AS $$
BEGIN
  update hints set cost = in_cost where exercise_id = in_exercise_id AND position  = in_position;
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
drop function add_task(UUID, int,text);
CREATE OR REPLACE FUNCTION add_task(in_exercise_id UUID, in_position int, in_content text) 
RETURNS void AS $$
DECLARE
max_position int;
BEGIN
  SET CONSTRAINTS ALL DEFERRED;
  select max(position) into max_position from tasks where exercise_id = in_exercise_id;
  if max_position > in_position then
    update tasks set position = position+1 where exercise_id = in_exercise_id AND position >= in_position;
    insert into tasks(exercise_id,position,content) values(in_exercise_id,in_position,in_content);
  else 
    insert into tasks(exercise_id,position,content) values(in_exercise_id,max_position+1,in_content);
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
  SET CONSTRAINTS ALL DEFERRED;
  update tasks set position = in_new_position where exercise_id = in_exercise_id AND position = in_old_position;
  update tasks set position = in_old_position where exercise_id = in_exercise_id AND position = in_new_position; 
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;



