

with recursive p(id,level,path,description) as (
  select id,0, '/'||id ,description from modules m left join module_parents mp on m.id = mp.child_id where mp.child_id is null

  union
  
  select m.id,level+1, path || '/' || m.id, m.description
  from modules m left join module_parents mp on m.id = mp.child_id inner join p on p.id = mp.parent_id

) select description, level from p order by level
