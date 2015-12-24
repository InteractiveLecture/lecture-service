\c lecture
--TODO unit test
drop function complete_exercise(UUID,UUID);
CREATE OR REPLACE FUNCTION complete_exercise(in_exercise_id UUID, in_user_id UUID) 
RETURNS void AS $$
DECLARE
beaten_exercises int;
exercise_module_id UUID;
exercise_topic_id UUID;
sum_points int;
BEGIN
  insert into exercise_progress_histories(user_id,exercise_id,amount,state,time) values(in_user_id,in_exercise_id,100,2,now());
  sum_points = 100;
  select count(*) into beaten_exercises from exercise_progress_histories where user_id = in_user_id AND exercise_id = in_exercise_id AND state = 2;

  select m.id,t.id
  into exercise_module_id ,exercise_topic_id 
  from exercises e
  inner join modules m on e.module_id = m.id
  inner join topics t on t.id = m.topic_id
  where e.id = in_exercise_id; 
  if  beaten_exercises > 2 then
    insert into module_progress_histories(user_id,module_id,amount,state,time) values(in_user_id,exercise_module_id,300,2,no());
    sum_points = sum_points + 300;
  end if;
  update topic_balances set amount = amount + sum_points where user_id = in_user_id AND topic_id = exercise_topic_id;
END;
$$ LANGUAGE plpgsql;

drop function complete_task(UUID,UUID);
CREATE OR REPLACE FUNCTION complete_task(in_task_id UUID, in_user_id UUID) 
RETURNS void AS $$
DECLARE
beaten_exercises int;
beaten_tasks int;
existend_tasks int;
var_exercise_id UUID;
exercise_module_id UUID;
exercise_topic_id UUID;
sum_points int;
BEGIN
  insert into task_completed_histories(user_id,task_id,time) values(in_user_id,in_task_id,now());
  select exercise_id into var_exercise_id from tasks where task_id = in_task_id;
  select max(position) into existend_tasks from tasks where exercise_id = var_exercise_id group by exercise_id;
  select count(tch.task_id) into beaten_tasks from tasks t inner join task_completed_histories tch on t.id = tch.task_id where tch.user_id = in_user_id AND t.exercise_id = var_exercise_id;

  if beaten_tasks == existend_tasks then
    return;
  end if;
  --if we are here, all tasks of the exercise have been beaten
  insert into exercise_progress_histories(user_id,exercise_id,amount,state,time) values(in_user_id,var_exercise_id,100,2,now());
  sum_points = 100;
  select count(*) into beaten_exercises from exercise_progress_histories where user_id = in_user_id AND exercise_id = in_exercise_id AND state = 2;

  select m.id,t.id
  into exercise_module_id ,exercise_topic_id 
  from exercises e
  inner join modules m on e.module_id = m.id
  inner join topics t on t.id = m.topic_id
  where e.id = in_exercise_id; 
  if  beaten_exercises > 2 then
    insert into module_progress_histories(user_id,module_id,amount,state,time) values(in_user_id,exercise_module_id,300,2,no());
    sum_points = sum_points + 300;
  end if;
  update topic_balances set amount = amount + sum_points where user_id = in_user_id AND topic_id = exercise_topic_id;
END;
$$ LANGUAGE plpgsql;

--TODO unit test
drop function purchase_hint(UUID,UUID);
CREATE OR REPLACE FUNCTION purchase_hint(in_hint_id UUID, in_user_id UUID) 
RETURNS int AS $$
DECLARE
hint_cost int;
user_balance int;
hint_topic_id UUID;
BEGIN
  if exists(select 1 from hint_purchase_histories where user_id = in_user_id AND hint_id = in_hint_id) then
    return 2;
  end if;

  if not exists(select 1 from hints where id = in_hint_id) then
    return 3;
  end if;
  select tb.amount, h.cost ,tb.topic_id into user_balance , hint_cost , hint_topic_id 
  from topic_balances tb 
  inner join modules m on m.topic_id = tb.topic_id
  inner join exercises e on e.module_id = m.id
  inner join tasks t on t.exercise_id = e.id
  inner join hints h on h.task_id = t.id
  where tb.user_id = in_user_id  AND h.id = in_hint_id;

  if user_balance < hint_cost then
    return 1;
  end if;

  insert into hint_purchase_histories(user_id,hint_id,amount,time) values(in_user_id,in_hint_id,hint_cost,now());

  update topic_balances set amount = amount - hint_cost where user_id = in_user_id AND topic_id = hint_topic_id;
  return 0;
END;
$$ LANGUAGE plpgsql;



