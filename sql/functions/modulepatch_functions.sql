\c lecture

--TODO unit test
drop function replace_module_description(UUID,text);
CREATE OR REPLACE FUNCTION replace_module_description(in_module_id UUID,in_new_description text) 
RETURNS void AS $$
  update modules set description = in_new_description where id = in_module_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_trees;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function remove_module_recommendation(UUID,UUID);
CREATE OR REPLACE FUNCTION remove_module_recommendation(in_module_id UUID,in_recommendation_id UUID) 
RETURNS void AS $$
  delete from module_recommendations where recommender_id = in_module_id AND recommended_id = in_recommendation_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function add_module_recommendation(UUID,UUID);
CREATE OR REPLACE FUNCTION add_module_recommendation(in_module_id UUID,in_recommendation_id UUID) 
RETURNS void AS $$
  insert into module_recommendations(recommender_id, recommended_id) values(in_module_id,in_recommendation_id);
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;



--TODO unit test
drop function add_module_script(UUID,UUID);
CREATE OR REPLACE FUNCTION add_module_script(in_module_id UUID,in_script_id UUID) 
RETURNS void AS $$
  update modules set script_id = in_script_id where id = in_module_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_trees;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;

--TODO unit test
drop function remove_module_script(UUID,UUID);
CREATE OR REPLACE FUNCTION remove_module_script(in_module_id UUID,in_video_id UUID) 
RETURNS void AS $$
  update modules set script_id = null where id = in_module_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_trees;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function add_module_video(UUID,UUID);
CREATE OR REPLACE FUNCTION add_module_video(in_module_id UUID,in_video_id UUID) 
RETURNS void AS $$
  update modules set video_id = in_video_id where id = in_module_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_trees;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;


--TODO unit test
drop function remove_module_video(UUID,UUID);
CREATE OR REPLACE FUNCTION remove_module_video(in_module_id UUID,in_video_id UUID) 
RETURNS void AS $$
  update modules set video_id = null where id = in_module_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_trees;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;

--TODO unit test
drop function add_exercise(UUID,UUID,varchar);
CREATE OR REPLACE FUNCTION add_exercise(in_exercise_id UUID, in_module_id UUID, in_backend varchar) 
RETURNS void AS $$
  insert into exercises(id,module_id,backend,version) values(in_exercise_id,in_module_id,in_backend,1);
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
$$ LANGUAGE sql;

--TODO unit test
drop function remove_exercise(UUID,UUID);
CREATE OR REPLACE FUNCTION remove_exercise(in_context_id UUID, in_exercise_id UUID) 
RETURNS void AS $$
BEGIN
  PERFORM check_module_context(in_context_id,in_exercise_id);
  delete from exercises where id = in_exercise_id;
  REFRESH MATERIALIZED VIEW CONCURRENTLY module_details;
END;
$$ LANGUAGE plpgsql;


drop function check_module_context(UUID,UUID);
CREATE OR REPLACE FUNCTION check_module_context(in_context_id UUID, in_exercise_id UUID) 
RETURNS void AS $$
BEGIN
  if NOT exists(select 1 from exercises e where e.id = in_exercise_id AND e.module_id= in_context_id) then
    RAISE EXCEPTION 'Operation out of scope.';
  end if;
END;
$$ LANGUAGE plpgsql;
