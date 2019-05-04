package neptune

var SQL_samples string = `SELECT 
                a.sample_age_ma, a.sample_depth_mbsf, c.water_depth, c.leg, c.site, c.hole, a.hole_id, 
                c.latitude, c.longitude, c.ocean_code, d.taxon_abundance, 
                TRIM(INITCAP(e.genus) || ' ' || LOWER(e.species) || ' ' || LOWER(COALESCE(e.subspecies,' '))) AS taxon, e.fossil_group 
            FROM 
                neptune_sample a, neptune_core b, neptune_hole_summary c, neptune_sample_taxa d, 
                neptune_taxonomy e  
            WHERE
                a.hole_id = b.hole_id 
                AND b.hole_id = c.hole_id 
                AND a.core = b.core 
                AND a.sample_id = d.sample_id 
                AND e.taxon_id = d.taxon_id            
            ORDER BY 
                a.sample_age_ma
	  LIMIT 1000`
