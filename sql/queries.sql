-- depth request
\echo 'depth request'
SELECT m.description, m.parents , (array_length(m.ancestors,1)-array_length(parents,1)) as depth 
FROM modules m 
WHERE 
  m.topic_id = 'b8c98f3e-bb7c-39e7-a3ce-e479c7892882' 
  AND  (array_length(m.ancestors,1)-array_length(parents,1))  < 4 
  AND  (array_length(m.ancestors,1)-array_length(parents,1)) >= 1 
ORDER BY depth;

-- tree request
\echo 'tree request'
SELECT m.description, m.parents
FROM modules m 
WHERE m.topic_id = 'b8c98f3e-bb7c-39e7-a3ce-e479c7892882' 
ORDER BY array_length(m.ancestors,1) NULLS FIRST;
