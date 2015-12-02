drop function if exists query_topics(integer,integer);
drop function if exists get_topic(UUID);
drop function if exists get_module_tree(UUID,int,int);

CREATE OR REPLACE FUNCTION query_topics(skip int, query_limit int)  returns json AS $$
select to_json(array_agg(row_to_json(o1)))
from(
  select t.id, t.name, t.description, t.version,(
    select array_agg(row_to_json(d)) 
    from ( select user_id,kind from topic_authority where topic_id = t.id)d ) as authorities 
  from topics t inner join topic_authority a on t.id = a.topic_id 
  LIMIT query_limit
  OFFSET skip
) o1;
$$ LANGUAGE sql;

CREATE OR REPLACE FUNCTION get_topic(in_topic_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1) into result
  from(
    select t.id, t.name, t.description, t.version,(
      select array_agg(row_to_json(d)) 
      from ( select user_id,kind from topic_authority where topic_id = in_topic_id)d ) as authorities 
    from topics t inner join topic_authority a on t.id = a.topic_id 
    where t.id = in_topic_id
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_module_tree(in_topic_id UUID, upper_bound int, lower_bound int)  returns json AS $$
DECLARE
result json;
BEGIN

  if lower_bound = -1  AND upper_bound = -1 then
    select to_json(array_agg(row_to_json(o1))) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id 
    ) o1;
    return result;
  end if;

  if lower_bound = -1 then 
    select to_json(array_agg(row_to_json(o1))) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level <= upper_bound 
    ) o1;

    return result;
  end if;

  if upper_bound = -1 then
    select to_json(array_agg(row_to_json(o1))) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level >= lower_bound
    ) o1;
    return result;
  end if;

  select to_json(array_agg(row_to_json(o1))) into result
  from(
    select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
    from module_trees m 
    where m.topic_id = in_topic_id AND m.level <= upper_bound AND m.level >= lower_bound
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;



DROP FUNCTION IF EXISTS get_module(UUID);
CREATE OR REPLACE FUNCTION get_module(in_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id,m.children,(
      select array_agg(row_to_json(exercises_aggregator)) from ( --aggregate exercises
        select ex.id, ex.backend, ex.version, (
          select array_agg(row_to_json(parts_aggregator)) from(
            select ta.content from tasks ta where ta.exercise_id = ex.id order by position
          ) parts_aggregator
          ) as parts, (
          select array_agg(row_to_json(hints_aggregator)) from (
            select hi.id from hints hi where hi.exercise_id = ex.id order by position
          ) hints_aggregator
        ) as hint_ids
        from exercises ex where ex.module_id = m.id
      ) exercises_aggregator
    ) as exercises
    from module_trees m 
    where m.id = in_id 
  ) o1;

  return result;
END;
$$ LANGUAGE plpgsql;



DROP FUNCTION IF EXISTS get_balances(UUID);
CREATE OR REPLACE FUNCTION get_balances(in_user_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select b.user_id, b.topic_id,b.amount from topic_balances b where b.user_id = in_user_id) o1;
  return result;
END;
$$ LANGUAGE plpgsql;




DROP FUNCTION IF EXISTS get_hint_purchase_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_hint_purchase_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select h.user_id, h.hint_id, h.amount, h.time 
    from hint_purchas_history h 
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;




DROP FUNCTION IF EXISTS get_module_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_module_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select h.user_id, h.module_id, h.reward, m.description, h.time
    from module_progress_histories h inner join modules m on h.module_id = m.id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS get_exercise_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_exercise_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select h.user_id, h.exercise_id, h.reward,h.time 
    from exercise_progress_histories h 
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS check_hint_purchase(UUID,UUID);
CREATE OR REPLACE FUNCTION check_hint_purchase(in_user_id UUID,in_hint_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select h.user_id, h.hint_id, h.amount, h.time 
    from hint_purchas_history h 
    where h.user_id = in_user_id AND h.hint_id = in_hint_id
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS get_hint(UUID,UUID);
CREATE OR REPLACE FUNCTION get_hint(in_user_id UUID,in_hint_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  if EXISTS( select 1 from hint_purchas_histories where user_id = in_user_id AND hint_id = in_hint_id) OR  
    EXISTS(select 1 
      from topic_authority ta inner join topics t on ta.topic_id = t.id 
      inner join modules m on t.id = m.topic_id
      inner join exercises e on e.module_id = m.id
      inner join hints h on h.exercise_id = e. id 
      where ta.user_id = in_user_id AND h.id = in_hint_id)
    then
    select row_to_json(o1)  into result from(
      select id, exercise_id, position, content, cost from hints where id = in_hint_id
    ) o1;
  end if;
  return result;
END;
$$ LANGUAGE plpgsql;


/*
DROP FUNCTION IF EXISTS get_current_modules_for_user(UUID);
CREATE OR REPLACE FUNCTION get_current_modules_for_user(in_user_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN

  select to_json(array_agg(row_to_json(o1))) into result from (
    select h1.module_id,m.description,t.id,t.description
    from module_progress_histories h1 inner join modules m on m.id = h1.module_id
    inner join topics t on t.id = m.topic_id
    where h1.state = 1 
    AND h1.user_id = in_user_id 
    AND NOT EXISTS(select 1 from module_progress_histories h2 where h1.module_id = h2.module_id and h2.state = 2)) o1;
  return result;
END;
$$ LANGUAGE plpgsql;
 */

DROP FUNCTION IF EXISTS get_next_modules_for_user(UUID);
CREATE OR REPLACE FUNCTION get_next_modules_for_user(in_user_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select to_json(array_agg(row_to_json(o1))) into result from (
    select mp.child_id as id,m.description,t.id as topic_id,t.name as topic_name
    from module_progress_histories mh 
    inner join module_parents mp on mp.parent_id = mh.module_id 
    inner join modules m on m.id = mp.child_id 
    inner join topics t on t.id = m.topic_id 
    where mp.child_id not in (select module_id from module_progress_histories where state = 2) 
    AND mh.state = 2 
    AND mh.user_id = in_user_id) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


