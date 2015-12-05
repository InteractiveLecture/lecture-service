\c lecture

drop materialized view if exists module_trees;


create materialized view module_trees as 
with recursive p(id,level,path,description,topic_id,video_id,script_id) as (
  select id,0, '/'||id ,description ,topic_id, video_id,script_id from modules m left join module_parents mp on m.id = mp.child_id where mp.child_id is null
  union
  select m.id,level+1, path || '/' || m.id, m.description, m.topic_id,m.video_id,m.script_id
  from modules m left join module_parents mp on m.id = mp.child_id inner join p on p.id = mp.parent_id
),
main as (
  select id, level, array_agg(path) as paths, description , topic_id ,video_id,script_id from p group by p.id,p.level,p.description, p.topic_id, p.video_id, p.script_id order by level),
children as (
  select mt.id,array_agg(p.child_id) as children from main mt left join module_parents p on mt.id = p.parent_id group by mt.id
) select m.id,m.level,m.paths,m.description,m.topic_id,m.video_id,m.script_id,c.children from main m inner join children c on m.id = c.id order by m.level;

ALTER MATERIALIZED VIEW module_trees OWNER TO lectureapp; 
create unique index module_trees_index on module_trees (id,level);

DROP FUNCTION IF EXISTS get_exercises_as_json(UUID);
CREATE OR REPLACE FUNCTION get_exercises_as_json(in_module_id UUID)  returns json AS $$
select json_agg(exercises_aggregator) from ( --aggregate exercises
  select ex.id, ex.backend, ex.version, (
    select json_agg(parts_aggregator) from(
      select ta.content from tasks ta where ta.exercise_id = ex.id order by position
    ) parts_aggregator
    ) as parts, (
    select json_agg(hints_aggregator) from (
      select hi.id from hints hi where hi.exercise_id = ex.id order by position
    ) hints_aggregator
  ) as hint_ids
  from exercises ex where ex.module_id = in_module_id
) exercises_aggregator;
$$ LANGUAGE sql;

DROP FUNCTION IF EXISTS get_recommendations_as_json(UUID);
CREATE OR REPLACE FUNCTION get_recommendations_as_json(in_module_id UUID)  returns json AS $$
select json_agg(recommendations_aggregator) from ( 
  select r.recommended_id as id, m.description as description,t.id as topic_id, t.name as topic_name
  from module_recommendations r 
  inner join modules m on m.id = r.recommended_id
  inner join topics t on m.topic_id = t.id
  where r.recommender_id = in_module_id
) recommendations_aggregator;
$$ LANGUAGE sql;


drop materialized view if exists module_details;
CREATE materialized view module_details AS 
  select o1.id as id, o1.level as level,row_to_json(o1) as details from(
    select 
    m.id,
    m.level,
    m.paths, 
    m.description,
    m.topic_id,
    m.video_id,
    m.script_id,
    m.children,
    get_exercises_as_json(m.id) as exercises,
    get_recommendations_as_json(m.id) as recommendations
    from module_trees m 
  ) o1;

ALTER MATERIALIZED VIEW module_details OWNER TO lectureapp; 
create unique index module_details_index on module_details(id,level);
