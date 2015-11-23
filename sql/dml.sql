
/*
original situation. We want to move D to be a new sibling of A
  A - B - C
        \
          D - E - F

*/



CREATE OR REPLACE FUNCTION move_module(module_id UUID, new_parent_id UUID) 
RETURNS void AS $$
DECLARE 
  old_parent UUID
BEGIN
  select into old_parent parent_id from module_parents where child_id = module_id;

  /*
situation after update
          E - F
        / 
  A - B - C
        \
          D
*/


  update module_parents set parent_id = old_parent where parent_id = module_id;

  /*
situation after delete 
      E-F
     / 
  A-B-C
   \
    D

*/
  update module_parents set parent_id = new_parent_id where child_id = module_id;

END;
$$ LANGUAGE plpgsql;

-- move node in tree
begin
  update modules 
  SET depth=1 
  where id=uuid_generate_v3(uuid_ns_url(),'module_5');
  delete from module_parents where child_id = uuid_generate_v3(uuid_ns_url(),'module_5');
end;
