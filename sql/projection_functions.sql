drop function if exists query_topics(integer,integer);
drop function if exists get_topic(UUID);
drop function if exists get_module_tree(UUID,int,int);

CREATE OR REPLACE FUNCTION query_topics(skip int, query_limit int)  returns table(jsonresult json) AS $$
select row_to_json(o1)
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


CREATE OR REPLACE FUNCTION get_module_tree(in_topic_id UUID, upper_bound int, lower_bound int)  returns table(jsonresult json) AS $$
BEGIN

  if lower_bound is null AND upper_bound is null then
    return QUERY select row_to_json(o1) 
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id 
    ) o1;
  end if;


  if lower_bound is null then 
    return QUERY select row_to_json(o1) 
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level >= lower_bound
    ) o1;
  end if;

  if upper_bound is null then
    return QUERY select row_to_json(o1) 
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level >= lower_bound
    ) o1;
  end if;

  return QUERY select row_to_json(o1) 
  from(
    select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
    from module_trees m 
    where m.topic_id = in_topic_id AND m.level <= upper_bound AND m.level >= lower_bound
  ) o1;
END;
$$ LANGUAGE plpgsql;



DROP FUNCTION IF EXISTS get_module(UUID);
CREATE OR REPLACE FUNCTION get_module(in_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id,(
      select array_agg(row_to_json(exercises_aggregator)) from ( --aggregate exercises
        select ex.id, ex.backend, ex.version, (
          select array_agg(row_to_json(parts_aggregator)) from(
              select ta.id, ta.task from tasks ta where ta.exercise_id = ex.id 
            ) parts_aggregator
          ) as parts, (
          select array_agg(row_to_json(hints_aggregator)) from (
            select hi.id from hints hi where hi.exercise_id = ex.id
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
    select h.user_id, h.hint_id, h.amount,
    from hint_purchas_history h 
    where h.user_id = in_user_id
    LIMIT in_limit
    OFFSET in_offset
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;




DROP FUNCTION IF EXISTS get_module_progress_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_module_progress_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
DECLARE
result json;
BEGIN
  select row_to_json(o1)  into result from(
    select h.user_id, h.module_id, h.reward, m.description
    from hint_purchas_history h  inner join modules m
    on h.module_id = m.module_id
    where b.user_id = in_user_id
    LIMIT in_limit
    OFFSET in_offset
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;
