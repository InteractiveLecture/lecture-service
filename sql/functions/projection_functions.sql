\c lecture
drop function if exists query_topics(integer,integer);
drop function if exists get_topic(UUID);
drop function if exists get_module_tree(UUID,int,int);

CREATE OR REPLACE FUNCTION get_authority_as_json(in_topic_id UUID)  returns json AS $$
select coalesce(json_agg(d),'[]') as authorities
from ( select user_id,kind from topic_authority where topic_id = in_topic_id)d;
$$ LANGUAGE sql;


CREATE OR REPLACE FUNCTION query_topics(skip int, query_limit int)  returns json AS $$
select json_agg(o1)
from(
  select t.id, t.name, t.description, t.version,(get_authority_as_json(t.id)) as authorities
  from topics t  
  LIMIT query_limit
  OFFSET skip
) o1;
$$ LANGUAGE sql;

CREATE OR REPLACE FUNCTION get_topic(in_topic_id UUID)  returns json AS $$
select row_to_json(o1) 
from(
  select t.id, t.name, t.description, t.version,(get_authority_as_json(in_topic_id)) as authorities
  from topics t
  where t.id = in_topic_id
) o1;
$$ LANGUAGE sql;


CREATE OR REPLACE FUNCTION get_module_tree(in_topic_id UUID, upper_bound int, lower_bound int)  returns json AS $$
DECLARE
result json;
BEGIN

  if lower_bound = -1  AND upper_bound = -1 then
    select json_agg(o1) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id 
    ) o1;
    return result;
  end if;

  if lower_bound = -1 then 
    select json_agg(o1) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level <= upper_bound 
    ) o1;

    return result;
  end if;

  if upper_bound = -1 then
    select json_agg(o1) into result
    from(
      select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
      from module_trees m 
      where m.topic_id = in_topic_id AND m.level >= lower_bound
    ) o1;
    return result;
  end if;

  select json_agg(o1) into result
  from(
    select m.id,m.level,m.paths, m.description,m.topic_id,m.video_id,m.script_id
    from module_trees m 
    where m.topic_id = in_topic_id AND m.level <= upper_bound AND m.level >= lower_bound
  ) o1;
  return result;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS get_balances(UUID);
CREATE OR REPLACE FUNCTION get_balances(in_user_id UUID)  returns json AS $$
select json_agg(o1) from(
  select b.user_id, b.topic_id,b.amount from topic_balances b where b.user_id = in_user_id) o1;
$$ LANGUAGE sql;

DROP FUNCTION IF EXISTS get_hint_purchase_history_base(UUID,int,int);
CREATE OR REPLACE FUNCTION get_hint_purchase_history_base(in_user_id UUID,in_limit int,in_skip int)  returns table(user_id UUID, hint_id UUID, amount smallint, event_time timestamp, exercise_id UUID)AS $$
BEGIN
  if in_limit = -1 AND in_skip = -1 then 
    return query select h.user_id, h.hint_id, h.amount, h.time,ta.exercise_id
    from hint_purchase_histories h inner join hints hi on hi.id = h.hint_id
    inner join tasks ta on ta.id = hi.task_id
    where h.user_id = in_user_id
    ORDER BY h.time;
  elsif in_limit = -1 then
    return query select h.user_id, h.hint_id, h.amount, h.time, ta.exercise_id 
    from hint_purchase_histories h inner join hints hi on hi.id = h.hint_id
    inner join tasks ta on ta.id = hi.task_id
    where h.user_id = in_user_id
    ORDER BY h.time
    OFFSET in_skip;

  elsif in_skip = -1 then
    return query select h.user_id, h.hint_id, h.amount, h.time, ta.exercise_id
    from hint_purchase_histories h inner join hints hi on hi.id = h.hint_id
    inner join tasks ta on ta.id = hi.task_id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit;
  else
    return query select h.user_id, h.hint_id, h.amount, h.time, ta.exercise_id
    from hint_purchase_histories h inner join hints hi on hi.id = h.hint_id
    inner join tasks ta on ta.id = hi.task_id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip;
  end if;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS get_hint_purchase_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_hint_purchase_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
select json_agg(o1) from(select * from get_hint_purchase_history_base(in_user_id,in_limit,in_skip)) o1;
$$ LANGUAGE sql;

DROP FUNCTION IF EXISTS get_hint_purchase_history(UUID,int,int,UUID);
CREATE OR REPLACE FUNCTION get_hint_purchase_history(in_user_id UUID,in_limit int,in_skip int, in_exercise_id UUID)  returns json AS $$
select json_agg(o1) from(select * from get_hint_purchase_history_base(in_user_id,in_limit,in_skip) where exercise_id = in_exercise_id) o1;
$$ LANGUAGE sql;


DROP FUNCTION IF EXISTS get_module_history_base(UUID,int,int);
CREATE OR REPLACE FUNCTION get_module_history_base(in_user_id UUID,in_limit int,in_skip int)  returns table(user_id UUID, module_id UUID, amount smallint, description text, event_time timestamp, event_type varchar, topic_id UUID)AS $$
BEGIN
  if in_limit = -1 AND in_skip = -1 then
    return query
    select h.user_id, h.module_id, h.amount, m.description, h.time, ps.description,m.topic_id
    from module_progress_histories h 
    inner join modules m on h.module_id = m.id
    inner join progress_state ps on h.state = ps.id
    where h.user_id = in_user_id
    ORDER BY h.time;

  elsif in_limit = -1 then
    return query
    select h.user_id, h.module_id, h.amount, m.description, h.time, ps.description,m.topic_id
    from module_progress_histories h 
    inner join modules m on h.module_id = m.id
    inner join progress_state ps on h.state = ps.id
    where h.user_id = in_user_id
    ORDER BY h.time
    OFFSET in_skip;
  elsif in_skip = -1 then
    return query
    select h.user_id, h.module_id, h.amount, m.description, h.time, ps.description, m.topic_id
    from module_progress_histories h 
    inner join modules m on h.module_id = m.id
    inner join progress_state ps on h.state = ps.id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit;
  else 
    return query
    select h.user_id, h.module_id, h.amount, m.description, h.time, ps.description,m.topic_id
    from module_progress_histories h 
    inner join modules m on h.module_id = m.id
    inner join progress_state ps on h.state = ps.id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip;
  end if;
END;
$$ LANGUAGE plpgsql;



DROP FUNCTION IF EXISTS get_module_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_module_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
select json_agg(o1) from(
  select * from get_module_history_base(in_user_id,in_limit,in_skip)
) o1;
$$ LANGUAGE sql;


DROP FUNCTION IF EXISTS get_module_history(UUID,int,int,UUID);
CREATE OR REPLACE FUNCTION get_module_history(in_user_id UUID,in_limit int,in_skip int, in_topic_id UUID)  returns json AS $$
select json_agg(o1) from(
  select * from get_module_history_base(in_user_id,in_limit,in_skip) where topic_id = in_topic_id
) o1;
$$ LANGUAGE sql;





DROP FUNCTION IF EXISTS get_exercise_base(UUID,int,int);
CREATE OR REPLACE FUNCTION get_exercise_base(in_user_id UUID,in_limit int,in_skip int)  returns table (user_id UUID, exercise_id UUID, amount smallint, event_time timestamp, description varchar, module_id UUID)AS $$
BEGIN
  if in_limit = -1 AND in_skip = -1 then
    return query
    select h.user_id, h.exercise_id, h.amount,h.time, ps.description, e.module_id
    from exercise_progress_histories h
    inner join progress_state ps on h.state = ps.id
    inner join exercises e on e.id = h.exercise_id
    where h.user_id = in_user_id
    ORDER BY h.time;
  elsif in_limit = -1 then
    return query
    select h.user_id, h.exercise_id, h.amount,h.time, ps.description, e.module_id
    from exercise_progress_histories h
    inner join progress_state ps on h.state = ps.id
    inner join exercises e on e.id = h. exercise_id
    where h.user_id = in_user_id
    ORDER BY h.time
    OFFSET in_skip;
  elsif in_skip = -1 then
    return query
    select h.user_id, h.exercise_id, h.amount,h.time, ps.description, e.module_id
    from exercise_progress_histories h
    inner join progress_state ps on h.state = ps.id
    inner join exercises e on e.id = h.exercise_id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit;

  else 
    return query
    select h.user_id, h.exercise_id, h.amount,h.time, ps.description, e.module_id
    from exercise_progress_histories h
    inner join progress_state ps on h.state = ps.id
    inner join exercises e on e.id = h.exercise_id
    where h.user_id = in_user_id
    ORDER BY h.time
    LIMIT in_limit
    OFFSET in_skip;
  end if;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION IF EXISTS get_exercise_history(UUID,int,int);
CREATE OR REPLACE FUNCTION get_exercise_history(in_user_id UUID,in_limit int,in_skip int)  returns json AS $$
select json_agg(o1) from (select * from get_exercise_base(in_user_id,in_limit,in_skip))o1;
$$ LANGUAGE sql;

DROP FUNCTION IF EXISTS get_exercise_history(UUID,int,int, UUID);
CREATE OR REPLACE FUNCTION get_exercise_history(in_user_id UUID,in_limit int,in_skip int, in_module_id UUID)  returns json AS $$
select json_agg(o1)  from(
  select * from get_exercise_base(in_user_id,in_limit,in_skip) where module_id = in_module_id
) o1;
$$ LANGUAGE sql;


DROP FUNCTION IF EXISTS get_hint(UUID,UUID);
CREATE OR REPLACE FUNCTION get_hint(in_user_id UUID,in_hint_id UUID)  returns json AS $$
DECLARE
result json;
BEGIN
  if EXISTS( select 1 from hint_purchase_histories where user_id = in_user_id AND hint_id = in_hint_id) OR  
    EXISTS(select 1 
      from topic_authority ta inner join topics t on ta.topic_id = t.id 
      inner join modules m on t.id = m.topic_id
      inner join exercises e on e.module_id = m.id
      inner join tasks tas  on tas.exercise_id = e.id
      inner join hints h on h.task_id = tas.id 
      where ta.user_id = in_user_id AND h.id = in_hint_id)
    then
    select row_to_json(o1)  into result from(
      select id, task_id, position, content, cost from hints where id = in_hint_id
    ) o1;
  end if;
  return result;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS get_next_modules_for_user(UUID);
CREATE OR REPLACE FUNCTION get_next_modules_for_user(in_user_id UUID)  returns json AS $$
select json_agg(o1) from (
  select mp.child_id as id,m.description,t.id as topic_id,t.name as topic_name
  from module_progress_histories mh 
  inner join module_parents mp on mp.parent_id = mh.module_id 
  inner join modules m on m.id = mp.child_id 
  inner join topics t on t.id = m.topic_id 
  where mp.child_id not in (select module_id from module_progress_histories where state = 2) 
  AND mh.state = 2 
  AND mh.user_id = in_user_id) o1;
$$ LANGUAGE sql;

DROP FUNCTION IF EXISTS get_one_exercise_as_json(UUID);
CREATE OR REPLACE FUNCTION get_one_exercise_as_json(in_exercise_id UUID)  returns json AS $$
select row_to_json(exercises_aggregator) from ( 
  select ex.id, ex.backend, ex.version, (get_tasks_as_json(ex.id)) as tasks
  from exercises ex where ex.id = in_exercise_id
) exercises_aggregator;
$$ LANGUAGE sql;



