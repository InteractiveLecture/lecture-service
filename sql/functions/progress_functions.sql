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
  insert into exercise_progress_histories(user_id,exercise_id,reward,state,time) values(in_user_id,in_exercise_id,100,2,now());
  sum_points = 100;
  select count(*) into beaten_exercises from exercise_progress_histories where user_id = in_user_id AND exercise_id = in_exercise_id AND state = 2;

  select m.module_id,t.topic_id 
  into exercise_module_id ,exercise_topic_id 
  from exercises e
  inner join modules m on e.module_id = m.id
  inner join topics t on t.id = m.topic_id
  where e.id = in_exercise_id and state = 2;
  if  beaten_exercises > 2 then
    insert into module_progress_histories(user_id,module_id,reward,state,time) values(in_user_id,exercise_module_id,300,2,no());
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
  select tb.balance, h.cost ,tb.topic_id into user_balance , hint_cost , hint_topic_id 
  from topic_balances tb 
  inner join modules m on m.topic_id = tb.topic_id
  inner join exercises e on e.module_id = m.id
  inner join hints h on h.exercise_id = e.id
  where tb.user_id = in_user_id  AND h.id = in_hint_id;

  if user_balance < hint_cost then
    return 1;
  end if;

  insert into hint_purchase_histories(user_id,hint_id,amount,time) values(in_user_id,in_hint_id,hint_cost,now());

  update topic_balances set balance = balance - hint_cost where user_id = in_user_id AND topic_id = hint_topic_id;
  return 0;
END;
$$ LANGUAGE plpgsql;



