

create user lectureapp;


create database lecture owner lectureapp;

\c lecture

drop table if exists topics cascade;
drop table if exists modules cascade ;
drop table if exists exercises cascade;
drop table if exists tasks cascade;
drop table if exists hints cascade;

drop table if exists topic_authority cascade;
drop table if exists topic_balances cascade ;
drop table if exists hint_purchase_histories cascade;
drop table if exists exercise_progress_histories cascade;
drop table if exists module_progress_histories cascade;
drop table if exists module_parents cascade;
drop table if exists module_recommendations cascade;
drop table if exists progress_state cascade;
drop table if exists task_completed_histories cascade;

create table topics (
  id UUID PRIMARY KEY,
  name VARCHAR(256) NOT NULL,
  description TEXT NOT NULL,
  version BIGINT NOT NULL CHECK(version > 0 )
);

create table modules(
  id UUID PRIMARY KEY,
  topic_id UUID REFERENCES topics(id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  video_id UUID,
  script_id UUID,
  version BIGINT NOT NULL CHECK(version >0)
);


create table exercises (
  id UUID PRIMARY KEY,
  module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
  backend varchar(256) NOT NULL,
  version BIGINT NOT NULL CHECK(version > 0)
);

create table tasks (
  id UUID PRIMARY KEY,
  exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
  position int NOT NULL CHECK(position > 0),
  content text NOT NULL,
  UNIQUE(exercise_id,position) DEFERRABLE INITIALLY DEFERRED

);

create table hints (
  id UUID PRIMARY KEY,
  task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
  position int NOT NULL CHECK(position > 0),
  content TEXT NOT NULL,
  cost INT NOT NULL CHECK(cost > 0),
  UNIQUE(task_id, position) DEFERRABLE INITIALLY DEFERRED
);

create table topic_authority (
  topic_id UUID references topics(id) ON DELETE CASCADE,
  user_id UUID,
  kind varchar(256)
);

create table topic_balances (
  user_id UUID,
  topic_id UUID references topics(id) ON DELETE CASCADE,
  amount INT NOT NULL CHECK(amount >= 0),
  PRIMARY KEY(user_id,topic_id)
);


create table hint_purchase_histories (
  user_id UUID,
  hint_id UUID references hints(id) on delete cascade,
  amount SMALLINT NOT NULL CHECK(amount >= 0),
  time timestamp,
  PRIMARY KEY(user_id,hint_id)
);


create table progress_state (
  id int primary key,
  description varchar NOT NULL
);

insert into progress_state(id,description) values(1,'BEGIN');
insert into progress_state(id,description) values(2,'FINISH');

create table module_progress_histories (
  user_id UUID,
  module_id UUID references modules(id) on delete cascade,
  amount SMALLINT NOT NULL,
  time timestamp,
  state int references progress_state(id),
  PRIMARY KEY (user_id,module_id,state)
);


create table exercise_progress_histories (
  user_id UUID,
  exercise_id UUID references exercises(id) on delete cascade,
  amount SMALLINT NOT NULL,
  time timestamp,
  state int references progress_state(id),
  PRIMARY KEY (user_id,exercise_id,state)
);

create table task_completed_histories(
  user_id UUID,
  task_id UUID references tasks(id) on delete cascade,
  time timestamp,
  PRIMARY KEY (user_id,task_id)
);

create table module_recommendations (
  recommender_id UUID REFERENCES modules(id) ON DELETE CASCADE,
  recommended_id UUID REFERENCES modules(id) ON DELETE CASCADE,
  PRIMARY KEY(recommended_id,recommender_id)
);

create table module_parents (
  child_id UUID REFERENCES modules(id) ON DELETE CASCADE,
  parent_id UUID REFERENCES modules(id) ON DELETE CASCADE,
  CONSTRAINT mp_pk PRIMARY KEY(child_id,parent_id) DEFERRABLE
);


