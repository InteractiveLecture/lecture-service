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
  insert into exercise_progress_histories(user_id,exercise_id,reward,time) values(in_user_id,in_exercise_id,100,now());
  sum_points = 100;
  select count(*) into beaten_exercises from exercise_progress_histories where user_id = in_user_id AND exercise_id = in_exercise_id;

  select m.module_id,t.topic_id 
  into exercise_module_id ,exercise_topic_id 
  from exercises e
  inner join modules m on e.module_id = m.id
  inner join topics t on t.id = m.topic_id
  where e.id = in_exercise_id;

  CASE 
    when beaten_exercises == 0 then
      insert into topic_balances(user_id,topic_id,amount) values(in_user_id,exercise_topic_id,sum_points);
    when beaten_exercises BETWEEN 1 AND 2 then 
      update topic_balances set amount = amount + sum_points where user_id = in_user_id AND topic_id = exercise_topic_id;
    when beaten_exercises > 2 then
      insert into module_progress_histories(user_id,module_id,reward,time) values(in_user_id,exercise_module_id,300,no());
      sum_points = sum_points + 300;
      update topic_balances set amount = amount + sum_points where user_id = in_user_id AND topic_id = exercise_topic_id;
  end case;

END;
$$ LANGUAGE plpgsql;



