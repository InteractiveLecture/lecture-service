CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

truncate table topics cascade;
truncate table modules cascade;
truncate table module_parents cascade;
truncate table module_recommendations cascade;
truncate table exercises cascade;
truncate table tasks cascade;
truncate table hints cascade;

insert into topics values ( uuid_generate_v3(uuid_ns_url(),'topic_1') ,'Grundlagen der Programmierung mit Java','bla', 1);
insert into topics values (uuid_generate_v3(uuid_ns_url(),'topic_2') ,'Descriptive Statistik','bla',1);



insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_1'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'foo',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);

insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_2'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bar',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_3'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bli',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_4'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bla',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_5'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'blubb',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);

insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_6'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bazz',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);

insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_7'),
  uuid_generate_v3(uuid_ns_url(),'topic_2'),
  'foobarbazz',
  uuid_generate_v3(uuid_ns_url(),'video_1'),
  uuid_generate_v3(uuid_ns_url(),'script_1'),
  1);



-- insert parents

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_2'),
  uuid_generate_v3(uuid_ns_url(),'module_1')
);

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_3'),
  uuid_generate_v3(uuid_ns_url(),'module_2')
);

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_4'),
  uuid_generate_v3(uuid_ns_url(),'module_3')
);

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_5'),
  uuid_generate_v3(uuid_ns_url(),'module_3')
);

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_6'),
  uuid_generate_v3(uuid_ns_url(),'module_5')
);

insert into module_parents values(
  uuid_generate_v3(uuid_ns_url(),'module_6'),
  uuid_generate_v3(uuid_ns_url(),'module_4')
);

-- insert recommendations

insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_1'),uuid_generate_v3(uuid_ns_url(),'module_7'));

insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_2'),uuid_generate_v3(uuid_ns_url(),'module_7'));
insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_3'),uuid_generate_v3(uuid_ns_url(),'module_7'));
insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_4'),uuid_generate_v3(uuid_ns_url(),'module_7'));

-- insert authority

insert into topic_authority values (uuid_generate_v3(uuid_ns_url(),'topic_1'),uuid_generate_v3(uuid_ns_url(),'user_1'),'OFFICER');

insert into topic_authority values (uuid_generate_v3(uuid_ns_url(),'topic_2'),uuid_generate_v3(uuid_ns_url(),'user_2'),'OFFICER');



-- insert exercises

insert into exercises values (uuid_generate_v3(uuid_ns_url(),'exercise_1'),uuid_generate_v3(uuid_ns_url(),'module_1'),'java',1);
insert into tasks values (uuid_generate_v3(uuid_ns_url(),'exercise_1'),1,'do something');

insert into tasks values (uuid_generate_v3(uuid_ns_url(),'exercise_1'),2,'do something');
insert into tasks values (uuid_generate_v3(uuid_ns_url(),'exercise_1'),3,'do something');


insert into hints values (uuid_generate_v3(uuid_ns_url(),'hint_1'),uuid_generate_v3(uuid_ns_url(),'exercise_1'),1,'take a hint',100);
insert into hints values (uuid_generate_v3(uuid_ns_url(),'hint_2'),uuid_generate_v3(uuid_ns_url(),'exercise_1'),2,'take another hint',100);



insert into exercises values (uuid_generate_v3(uuid_ns_url(),'exercise_2'),uuid_generate_v3(uuid_ns_url(),'module_1'),'java',1);
insert into exercises values (uuid_generate_v3(uuid_ns_url(),'exercise_3'),uuid_generate_v3(uuid_ns_url(),'module_1'),'java',1);

