create user lectureapp;

create database lecture owner lectureapp;

\c lecture

drop table if exists topics cascade;
drop table if exists modules cascade ;
drop table if exists exercises cascade;
drop table if exists hints cascade;

drop table if exists topic_authority cascade;
drop table if exists topic_balances cascade ;
drop table if exists hint_purchase_histories cascade;
drop table if exists exercise_progress_histories cascade;
drop table if exists module_progress_histories cascade;
drop table if exists module_parents cascade;
drop table if exists module_recommendations cascade;




create table topics (
  id UUID PRIMARY KEY,
  name VARCHAR(256),
  description TEXT
);

create table modules(
  id UUID PRIMARY KEY,
  topic_id UUID REFERENCES topics(id),
  description TEXT NOT NULL,
  last_modified BIGINT NOT NULL CHECK(last_modified >0)
);


create table exercises (
  id UUID PRIMARY KEY,
  module_id UUID REFERENCES modules(id),
  task text NOT NULL,
  backend varchar(256) NOT NULL,
  last_modified BIGINT NOT NULL CHECK(last_modified > 0)
);

create table hints (
  id UUID PRIMARY KEY,
  exercise_id UUID REFERENCES exercises(id),
  content TEXT NOT NULL,
  cost SMALLINT NOT NULL CHECK(cost > 0)
);



create table topic_authority (
  lecture_id UUID references topics(id),
  user_id UUID,
  kind varchar(256)
);

create table topic_balances (
  user_id UUID,
  topic_id UUID references topics(id),
  amount SMALLINT NOT NULL CHECK(amount >= 0),
  PRIMARY KEY(user_id,topic_id)
);


create table hint_purchase_histories (
  user_id UUID,
  hint_id UUID references hints(id),
  amount SMALLINT NOT NULL CHECK(amount >= 0),
  PRIMARY KEY(user_id,hint_id)
);


create table module_progress_histories (
  user_id UUID,
  module_id UUID references modules(id),
  reward SMALLINT NOT NULL CHECK(reward > 0),
  PRIMARY KEY (user_id,module_id)
);


create table exercise_progress_histories (
  user_id UUID,
  exercise_id UUID references exercises(id),
  reward SMALLINT NOT NULL CHECK(reward > 0 ),
  PRIMARY KEY (user_id,exercise_id)
);

create table module_recommendations (
  recommender_id UUID REFERENCES modules(id),
  recommended_id UUID REFERENCES modules(id),
  PRIMARY KEY(recommended_id,recommender_id)
);

create table module_parents (
  child_id UUID REFERENCES modules(id),
  parent_id UUID REFERENCES modules(id),
  CONSTRAINT mp_pk PRIMARY KEY(child_id,parent_id) DEFERRABLE
);


