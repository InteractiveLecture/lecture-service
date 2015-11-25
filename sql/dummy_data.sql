CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

truncate table topics cascade;
truncate table modules cascade;
truncate table module_parents cascade;
truncate table module_recommendations cascade;

insert into topics values ( uuid_generate_v3(uuid_ns_url(),'topic_1') ,'Grundlagen der Programmierung mit Java','bla', 1);
insert into topics values (uuid_generate_v3(uuid_ns_url(),'topic_2') ,'Descriptive Statistik','bla',1);



insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_1'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'foo',
  1);

insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_2'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bar',
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_3'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bli',
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_4'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bla',
  1);
insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_5'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'blubb',
  1);

insert into modules values (
  uuid_generate_v3(uuid_ns_url(),'module_6'),
  uuid_generate_v3(uuid_ns_url(),'topic_1'),
  'bazz',
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

--insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_1'),uuid_generate_v3(uuid_ns_url(),'module_7'));

--insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_2'),uuid_generate_v3(uuid_ns_url(),'module_7'));
--insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_3'),uuid_generate_v3(uuid_ns_url(),'module_7'));
--insert into module_recommendations values (uuid_generate_v3(uuid_ns_url(),'module_4'),uuid_generate_v3(uuid_ns_url(),'module_7'));
