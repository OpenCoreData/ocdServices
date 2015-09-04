package janus

func SQLCalls() {
	sql_chemcarb := ` select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  s.top_interval*100.0, s.bottom_interval*100.0,
  get_depths(x.section_id,'STD',s.top_interval,0,0) mbsf,
  null,
  avg(decode(cca.analysis_code,'INOR_C',cca.analysis_result)),
  avg(decode(cca.analysis_code,'CaCO3', cca.analysis_result)),
  avg(decode(cca.analysis_code,'TOT_C', cca.analysis_result)),
  avg(decode(cca.analysis_code,'ORG_C', cca.analysis_result)),
  avg(decode(cca.analysis_code,'NIT',   cca.analysis_result)),
  avg(decode(cca.analysis_code,'SUL',   cca.analysis_result)),
  avg(decode(cca.analysis_code,'H',     cca.analysis_result))
from hole h, section x, sample s,
  chem_carb_sample ccs, chem_carb_analysis cca
where h.leg = x.leg and h.site = x.site and h.hole = x.hole and
  x.section_id = s.sam_section_id and
  s.sample_id = ccs.sample_id and
  s.location = ccs.location and
  ccs.run_id = cca.run_id
    {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
   group by x.leg, x.site, x.hole, x.core, x.core_type, x.section_type, x.section_number, s.top_interval, s.bottom_interval, x.section_id
   order by x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, s.top_interval
    `

	// --   if (legv != null) leg = " and x.leg = ${legv} "
	// -- if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// -- if (holev != null) hole = " and x.hole = upper('${holev}') "

	// -- def sqlfinal = sql + leg + site + hole + " group by x.leg, x.site, x.hole, x.core, x.core_type, x.section_type, x.section_number, s.top_interval, s.bottom_interval, x.section_id\n" +
	// --         "order by x.leg, x.site, x.hole, x.core,\n" +
	// --         "  x.core_type, x.section_number, s.top_interval"

	sql_roundSumTotalCore := ` select c.leg
     , c.site
     , c.hole
     , c.core
     , c.core_type
     , to_char(c.time_on_deck + c.offset/24, 'mm/dd/rr hh24mi')
     , c.top_depth
     , round(c.advancement*decode(instr('ABCDEHMNPRVXYZ',c.core_type), 0, 0, 1),2)
     , sum_liner_lengths(c.leg,c.site,c.hole,c.core,c.core_type)
     , sum_curated_lengths(c.leg,c.site,c.hole,c.core,c.core_type)
     , trim(get_core_comments(c.leg, c.site, c.hole, c.core, c.core_type))
     , get_min_core_depth(c.leg, c.site, c.hole, c.core, 'STD', 0, 0)
     , c.meter_comp_depth
     , c.advancement
     , c.drilling_time
    from
     core c
    where 1=1
    {{if .Leg }} and c.leg = {{.Leg}} {{end}}
    {{if .Site }} and c.site = {{.Site}} {{end}}
    {{if .Hole }} and c.hole = upper('{{.Hole}}') {{end}}
    and instr('ABCDEHMNPRVXYZ',c.core_type) > 0 order by c.leg, c.site, c.hole, c.core, c.core_type
    `

	// -- if (legv != null) leg = " and c.leg = ${legv} "
	// --  if (sitev != null) site = " and c.site = ${sitev} " ?: " "
	// --  if (holev != null) hole = " and c.hole = upper('${holev}') "
	// --  def sqlfinal = sql + leg + site + hole + " and instr('ABCDEHMNPRVXYZ',c.core_type) > 0 order by c.leg, c.site, c.hole, c.core, c.core_type"

	sql_sampleCount := ` select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  s.top_interval * 100.0, s.bottom_interval * 100.0,
  get_depths(x.section_id,'STD',s.top_interval,0,0),
  null,
  sr.request_number||sr.request_part_no,
  s.volume, s.piece||s.sub_piece,
  s.sam_archive_working,
  s.s_c_sampling_code, reqc.catwalk_sample,
  s.sam_sample_code_lab, labc.catwalk_sample,
  s.sample_id, s.location, s.sam_repository,
  to_char(s.timestamp,'mm/dd/yyyy hh24:mi'), s.sample_comment
from sample s, sample_code reqc, sample_request sr, sample_code_lab labc, section x
where (s.sam_section_id = x.section_id)
  and (s.s_c_leg = reqc.leg(+)
       and s.s_c_sampling_code = reqc.sampling_code(+))
  and (reqc.code_samrqst_id = sr.samrqst_id(+))
  and (s.sam_sample_code_lab = labc.sample_code_lab(+))
 {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
 order by 1,2,3,4,6,8
`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --        if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --        if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --        def sqlfinal = sql + leg + site + hole + " order by 1,2,3,4,6,8"

	sql_MSL_SECTION_COUNT := `select x.leg, x.site, x.hole, x.core, x.core_type,
     nvl(decode(section_type,'C','CC ',to_char(x.section_number,'90')), ' '),
     mst_top_interval * 100.0,
     -1.0,
     nvl(get_depths(x.section_id,'STD', mxd.mst_top_interval,0,0), -1) mbsf,
     Meas_Susceptibility_Mean, decode(bkgd_elapsed_zero_time,0,null,
       Meas_Susceptibility_Mean - (Bkgd_Susceptibility *
       sample_elapsed_zero_time / bkgd_elapsed_zero_time))
    from hole h, section x, msl_section mx, msl_section_data mxd
   where x.leg = h.leg and x.site = h.site and
              x.hole = h.hole and
         x.section_id = mx.section_id and
         mx.msl_id = mxd.msl_id
 {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
order by x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, mxd.mst_top_interval
   `

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --      if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --      if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --      def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, mxd.mst_top_interval"

	sql_NGR_SECTION_COUNT := ` select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  nxd.mst_top_interval*100.0,
  get_depths(x.section_id,'STD', nxd.mst_top_interval, 0, 0) mbsf,
  null,
  nxd.total_counts_sec - nb.total_counts_sec,
  nxd.total_counts_sec
 from section x, ngr_section nx, ngr_section_data nxd, ngr_background nb
where x.section_id = nx.section_id
  and nx.ngr_id = nxd.ngr_id
  and nx.energy_background_id = nb.energy_background_id(+)
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
order by x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, nxd.mst_top_interval`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --        if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --        if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --        def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, nxd.mst_top_interval"

	sql_PWL_SECTION_COUNT := ` select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  psd.mst_top_interval * 100,
  get_depths(x.section_id,'STD', psd.mst_top_interval,0,0) depthMbsf,
  null,
  nvl(psd.pwl_velocity,
    janus.PWL_Velocity(ps.pwl_id, psd.mst_top_interval,
      'VELOCITY', ps.liner_standard_id)
  )
from section x, pwl_section ps, pwl_section_data psd
where x.section_id = ps.section_id
  and ps.pwl_id = psd.pwl_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
order by x.leg, x.site, x.hole, x.core,  x.core_type, x.section_number, psd.mst_top_interval
`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --        if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --        if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --        def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core,  x.core_type, x.section_number, psd.mst_top_interval"

	sql_PWS_SECTION_COUNT := ` select
 x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , pxd.pp_top_interval * 100.0
 , pxd.pp_bottom_interval * 100.0
 , get_depths(x.section_id, 'STD', pxd.pp_top_interval, 0, 0)
 , null
 , px.direction
 , nvl(pxd.pws3_velocity, janus.pws3_velocity(px.pws_id, pxd.pp_top_interval, px.direction))
from
 pws3_section px
 , pws3_section_data pxd
 , section x
 , pws3_calibration pc
where
 pxd.pws_id = px.pws_id
 and px.section_id = x.section_id
 and px.pws_calibration_id = pc.pws_calibration_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
order by x.leg, x.site, x.hole, x.core, x.section_number, pxd.pp_top_interval
`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --         if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --         def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core, x.section_number, pxd.pp_top_interval"

	sql_MAD_SAMPLE_COUNT := ` (select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  s.top_interval * 100.0, s.bottom_interval * 100.0,
  get_depths(x.section_id,'STD', s.top_interval,0,0) depthMbsf,
  null,
  water_content_wet(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  water_content_dry(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  bulk_density(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  dry_density(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  grain_density(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  porosity(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  void_ratio(substr(msd.method,1,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  substr(msd.method,1,1) calcMethod,
  msd.method,
  msd.comments,
  to_char(msd.beaker_date_time, 'DD-MON-YYYY HH24:MI'),
  msd.water_content_bulk,
  msd.water_content_solids,
  msd.bulk_density,
  msd.dry_density,
  msd.grain_density,
  msd.porosity,
  msd.void_ratio
from
  section x, sample s,
  mad_sample_data msd, mad_beaker_history mbh
where
  substr(msd.method,1,1) is not null
  and x.section_id = s.sam_section_id
  and (s.sample_id,s.location) in ((msd.sample_id,msd.location))
  and msd.mad_beaker_id = mbh.mad_beaker_id(+)
  and msd.beaker_date_time = mbh.beaker_date_time(+)
 and x.leg = 210
 and x.site = 1277
union
select x.leg, x.site, x.hole, x.core, x.core_type,
  x.section_number, x.section_type,
  s.top_interval * 100.0, s.bottom_interval * 100.0,
  get_depths(x.section_id,'STD', s.top_interval,0,0) depthMbsf,
  null,
  water_content_wet(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  water_content_dry(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  bulk_density(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  dry_density(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  grain_density(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  porosity(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  void_ratio(substr(msd.method,2,1), msd.mass_wet_and_beaker, msd.mass_dry_and_beaker, mbh.beaker_mass, msd.vol_wet_and_beaker, msd.vol_dry_and_beaker, mbh.beaker_volume),
  substr(msd.method,2,1) calcMethod,
  msd.method,
  msd.comments,
  to_char(msd.beaker_date_time, 'DD-MON-YYYY HH24:MI'),
  msd.water_content_bulk,
  msd.water_content_solids,
  msd.bulk_density,
  msd.dry_density,
  msd.grain_density,
  msd.porosity,
  msd.void_ratio
from
  section x, sample s,
  mad_sample_data msd, mad_beaker_history mbh
where
  substr(msd.method,2,1) is not null
  and x.section_id = s.sam_section_id
  and (s.sample_id,s.location) in ((msd.sample_id,msd.location))
  and msd.mad_beaker_id = mbh.mad_beaker_id(+)
  and msd.beaker_date_time = mbh.beaker_date_time(+)
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
) order by 1, 2, 3, 4, 5, 6, 8, calcMethod 

`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --         if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --         def sqlfinal = sql + leg + site + hole + ") order by 1, 2, 3, 4, 5, 6, 8, calcMethod "

	sql_THERMCON_COUNT = ` select
 x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , td.pp_top_interval * 100.0
 , td.pp_bottom_interval * 100.0
 , get_depths(x.section_id,'STD', td.pp_top_interval,0,0) depthMbsf
, null depthOther
 , td.tcon_comment
 , td.tcon_probe_half_full
 , td.tcon_proc_thermcon
 , td.system_id
 , td.tcon_probe_num
from
 section x
 , tcon_data td
where x.section_id = td.section_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
order by x.leg, x.site, x.hole, x.core, x.section_number,  td.pp_top_interval
`

	// -- if (legv != null) leg = " and x.leg = ${legv} "
	// --        if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	// --        if (holev != null) hole = " and x.hole = upper('${holev}') "

	// --        def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core, x.section_number,  td.pp_top_interval"

	sql_SHEAR_STRENGTH_COUNT := ` select
 x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , ad.pp_top_interval * 100.0
 , get_depths(x.section_id,'STD', ad.pp_top_interval, 0, 0) depthMbsf
 , null
 , ad.strength_reading
 , to_char(xd.run_date_time)
 , xd.direction
 , xd.range
 , ad.comments
from
 section x
 , tor_section_data xd
 , tor_sample_data ad
where
 x.section_id = xd.section_id
 and xd.tor_id = ad.tor_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
    order by x.leg , x.site , x.hole , x.core , x.core_type , x.section_number , ad.pp_top_interval
`

	//   if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by x.leg , x.site , x.hole , x.core , x.core_type , x.section_number , ad.pp_top_interval"

	sql_MS2F_SECTION_COUNT := ` select x.leg, x.site, x.hole, x.core, x.core_type,
  nvl(decode(x.section_type,'C','CC ',to_char(x.section_number,'90')), ' '),
  ms2f_top_interval * 100.0,
  -1.0,
  nvl(get_depths(x.section_id,'STD', mxd.ms2f_top_interval,0,0), -1) mbsf,
  meas_susceptibility_mean, decode(drift_corr_susceptibility,0,null,
    meas_susceptibility_mean - average_mag_susc)
from hole h, section x, ms2f_section mx, ms2f_section_data mxd
where x.leg = h.leg and x.site = h.site and
      x.hole = h.hole and
      x.section_id = mx.section_id and
      mx.ms2f_id = mxd.ms2f_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
    order by x.leg, x.site, x.hole, x.core,  x.core_type, x.section_number, mxd.ms2f_top_interval
`

	//  if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site, x.hole, x.core,  x.core_type, x.section_number, mxd.ms2f_top_interval"

	sql_DHT_APCT_RUN_COUNT := ` select
  dar.apct_leg,
  dar.apct_site,
  dar.apct_hole,
  dar.apct_core,
  dar.apct_core_type,
  c.top_depth,
  null,
 c.top_depth + c.advancement,
  null,
  dar.apct_depth_comment,
  datr.apct_best_fit_temp_t0,
  datr.apct_best_fit_error_rms,
  datr.apct_mudline_temp,
  dac.apct_tool_name,
  dar.apct_run_comment
from
  core c,
  dht_apct_run dar,
  dht_apct_calib dac,
  dht_apct_tfit_results datr
where
  (c.leg, c.site, c.hole, c.core, c.core_type)
    in ((dar.apct_leg, dar.apct_site, dar.apct_hole, dar.apct_core, dar.apct_core_type))
  and (dac.apct_tool_id, dac.apct_calib_date_time)
  in (( dar.apct_tool_id, dar.apct_calib_date_time))
  and dar.apct_run_id = datr.apct_run_id(+)
{{if .Leg }} and dar.apct_leg = {{.Leg}} {{end}}
    {{if .Site }} and dar.apct_site = {{.Site}} {{end}}
    {{if .Hole }} and dar.apct_hole = upper('{{.Hole}}') {{end}}
order by 1, 2, 3, 4, 5, 6, 7
`

	// if (legv != null) leg = " and dar.apct_leg = ${legv} "
	//         if (sitev != null) site = " and dar.apct_site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and dar.apct_hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by 1, 2, 3, 4, 5, 6, 7"

	sql_SPLICER_COUNT := ` select
 lb.leg
 , lb.site
 , lb.sort_key
 , lb.hole_end
 , ub.mcd_end
 , lb.mcd_start
from
 splice ub, splice lb
where
 lb.sort_key = ub.sort_key(+) + 1
 and lb.leg = ub.leg(+)
 and lb.site = ub.site(+)
 and lb.hole_start = ub.hole_end(+)
 
 {{if .Leg }} and lb.leg = {{.Leg}} {{end}}
    {{if .Site }} and b.site = {{.Site}} {{end}}
    {{if .Hole }} and lb.hole_end = upper('{{.Hole}}') {{end}}
 order by lb.leg, lb.site, lb.sort_key

`

	//  if (legv != null) leg = " and lb.leg = ${legv} "
	//         if (sitev != null) site = " and lb.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and lb.hole_end = upper('${holev}') "   // todo check if this edit of adding _end to hole is propper

	//         def sqlfinal = sql + leg + site + hole + "order by lb.leg, lb.site, lb.sort_key"

	sql_TENSOR_CORE_COUNT := ` select rez.leg, rez.site,
  rez.hole, rez.core,
  rez.core_type,
  rez.hole_azimuth, rez.hole_inclination,
  rez.reorientation_angle_motf,
  rez.reorientation_angle_mtf,
  run.site_variation, c.top_depth
from tensor_tool_results rez, tensor_tool_runs run, hole h, core c
where c.leg = h.leg and c.site = h.site and
  c.hole = h.hole and rez.leg = c.leg and
  rez.site = c.site and
  rez.hole = c.hole and
  rez.core = c.core and
  rez.core_type = c.core_type and
  run.leg(+) = rez.leg and
  run.site(+) = rez.site and
  run.hole(+) = rez.hole and
  run.start_core(+) = rez.core
  {{if .Leg }} and rez.leg = {{.Leg}} {{end}}
    {{if .Site }} and rez.site = {{.Site}} {{end}}
    {{if .Hole }} and rez.hole_end = upper('{{.Hole}}') {{end}}
`

	//  if (legv != null) leg = " and rez.leg = ${legv} "
	//         if (sitev != null) site = " and rez.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and rez.hole = upper('${holev}') "

	sql_PALEO_SAMPLE_COUNT := ` select
 fg.fossil_group,
 fg.fossil_group_name,
 sci.scientist_id,
 s.sample_id, s.location,
 trim(sci.last_name),
 trim(sci.first_name),
 x.leg, x.site, x.hole, x.core, x.core_type,
 x.section_number, x.section_type,
 s.top_interval*100.0,
 get_depths(x.section_id, 'STD', s.top_interval, 0, 0) depthMbsf,
 null depthOther,
 trim(gac_old.geologic_age_name),
 trim(gac_yng.geologic_age_name),
 trim(zone_bottom.zone_abbrev),
 trim(zone_top.zone_abbrev),
 trim(fga.group_abundance_name),
 trim(p.preservation_name),
 replace(ps.paleo_sample_comment, chr(10), ' ')
from paleo_sample ps, fossil_group fg,
 fossil_group_abundance fga, sample s,
 section x, hole h, preservation p, scientist sci,
 geologic_age_concept gac_yng,
 geologic_age_concept gac_old,
 zone_concept zone_top,
 zone_concept zone_bottom
where (x.leg, x.site, x.hole) in ((h.leg, h.site, h.hole))
 and s.sam_section_id = x.section_id
 and (ps.sample_id, ps.location) in ((s.sample_id, s.location))
 and ps.scientist_id=sci.scientist_id
 and ps.fossil_group = fga.fossil_group(+)
 and ps.sample_group_abundance = fga.group_abundance(+)
 and ps.fossil_group = p.fossil_group(+)
 and ps.sample_preservation = p.preservation(+)
 and ps.fossil_group = fg.fossil_group
 and ps.geologic_age_young = gac_yng.geologic_age_id(+)
 and ps.geologic_age_old = gac_old.geologic_age_id(+)
 and ps.zone_old = zone_bottom.zone_id(+)
 and ps.zone_young = zone_top.zone_id(+)
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}

order by fg.fossil_group_name, sci.last_name, x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, s.top_interval
`

	//  if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by fg.fossil_group_name, sci.last_name, x.leg, x.site, x.hole, x.core, x.core_type, x.section_number, s.top_interval"

	sql_AGEPROFILE_COUNT := ` select
 ap.leg
 , ap.site
 , ap.hole
 , nvl(ap.datum_fossil_group, 0)
 , nvl(ap.datum_depth, 0.0)
 , nvl(ap.datum_depth, 0.0) + nvl(ap.datum_depth_error, 0.0)
 , nvl(ap.datum_age, 0.0)
 , nvl(ap.datum_age, 0.0) + nvl(ap.datum_age_error, 0.0)
 , nvl(ap.ageprofile_datum_id, 0)
 , nvl(ap.ageprofile_datum_description, ' ')
 , ap.ageprofile_datum_type
 , tc.GENUS_SUBGENUS
 , tc.SPECIES_SUBSPECIES
 , nvl(ap.depth_map_type, ' ')
 , nvl(ap.mcd_flag, ' ')
 , nvl(ap.compression_flag, ' ')
from
 ageprofile ap
 , taxon_concept tc
 , datum_concept dc
where
 ap.ageprofile_datum_id = dc.datum_id
 and ap.ageprofile_taxon_id = tc.taxon_id

{{if .Leg }} and leg = {{.Leg}} {{end}}
    {{if .Site }} and site = {{.Site}} {{end}}
    {{if .Hole }} and hole = upper('{{.Hole}}') {{end}}

order by ap.leg, ap.site, ap.hole, ap.datum_depth
`

	//  if (legv != null) leg = " and leg = ${legv} "
	//         if (sitev != null) site = " and site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and hole = upper('${holev}') "
	// def sqlfinal = sql + leg + site + hole + "order by ap.leg, ap.site, ap.hole, ap.datum_depth"

	sql_AGE_DATAPOINT_COUNT := ` select
 age_model_control_pts.leg, age_model_control_pts.site,
 age_model_control_pts.hole, age_model_type,
 depth, nvl(age, -1), nvl(age_model_control_pt_comment, ' ')
from
 hole h, age_model_control_pts
where h.leg = age_model_control_pts.leg
 and h.site = age_model_control_pts.site
 and h.hole = age_model_control_pts.hole

{{if .Leg }} and age_model_control_pts.leg = {{.Leg}} {{end}}
    {{if .Site }} and age_model_control_pts.site = {{.Site}} {{end}}
    {{if .Hole }} and age_model_control_pts.hole = upper('{{.Hole}}') {{end}}
order by age_model_control_pts.leg, age_model_control_pts.site, age_model_control_pts.hole, age_model_type, depth
`

	// if (legv != null) leg = " and age_model_control_pts.leg = ${legv} "
	//         if (sitev != null) site = " and age_model_control_pts.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and age_model_control_pts.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by age_model_control_pts.leg, age_model_control_pts.site, age_model_control_pts.hole, age_model_type, depth"

	sql_XRD_IMAGE_COUNT := ` select
 x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , get_depths(pdi.section_id, 'STD', 0, 0, 0)
 , get_depths(pdi.section_id, 'STD', 0, 0, 0)
 , dt.data_type_abbr
 , pdi.url
 , pdi.top_interval
 , fg.fossil_group_name
 , pdi.page_id
from
 prime_data_image pdi
 , section x
 , data_type dt
 , fossil_group fg
where
 pdi.section_id = x.section_id
 and pdi.data_type_id = dt.data_type_id
 and pdi.fossil_group = fg.fossil_group(+)
 {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}

`

	//  if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	// def sqlfinal = sql + leg + site + hole

	sql_XRF_SAMPLE_COUNT := ` SELECT xsa.sample_id,
  x.leg, x.site, x.hole, x.core, x.core_type,
  nvl(decode(section_type,'C','CC ', to_char(x.section_number,'90')), ' '),
  top_interval*100, bottom_interval * 100,
  -1.0,
  nvl(get_depths(x.section_id, 'STD', s.top_interval,0,0),-1),
  nvl(xsa.xrf_run_identifier, ' '),
  nvl(xsa.xrf_replicate, ' '),
  nvl(xsa.xrf_analysis_code, ' '),
  nvl(xsa.xrf_analysis_result, 0),
  nvl(xsa.analysis_units, ' '),
  xs.xrf_analysis_type,
  nvl(xst.sample_type, ' '),
  xs.bead_loi,
  nvl(xs.xrf_comment, ' ')
FROM xrf_sample_analysis xsa, sample s, section x, hole h, xrf_sample xs, xrf_sample_type xst
WHERE (x.leg, x.site, x.hole) in ((h.leg, h.site, h.hole)) and
  x.section_id = s.sam_section_id and
  s.location = xs.location and
  s.sample_id = xs.sample_id and
  xs.location = xsa.location and
  xs.sample_id = xsa.sample_id and
  xs.xrf_replicate = xsa.xrf_replicate and
  xs.xrf_run_identifier = xsa.xrf_run_identifier and
  xs.leg = xsa.leg and
  xs.xrf_analysis_type = xsa.xrf_analysis_type and
  xs.sample_type_id = xst.sample_type_id(+) and
  xs.system_id = 31
  {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
  order by x.leg, x.site,  x.hole, x.core, x.section_number,   s.top_interval, xsa.sample_id, xsa.xrf_run_identifier, xsa.xrf_replicate

`

	// if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site,  x.hole, x.core, x.section_number,   s.top_interval, xsa.sample_id, xsa.xrf_run_identifier, xsa.xrf_replicate"

	sql_ICP_SAMPLE_COUNT := ` SELECT xsa.sample_id,
  x.leg, x.site, x.hole, x.core, x.core_type,
  nvl(decode(x.section_type,'C','CC ', to_char(x.section_number,'90')), ' '),
  s.top_interval*100.0, s.bottom_interval * 100.0,
  nvl(get_depths(x.section_id, 'STD', s.top_interval,0,0),-1),
  -1.0,
  nvl(xsa.xrf_run_identifier, ' '),
  nvl(xsa.xrf_replicate, ' '),
  nvl(xsa.xrf_analysis_code, ' '),
  nvl(xsa.xrf_analysis_result,0),
  nvl(xsa.analysis_units, ' '),
  xs.xrf_analysis_type,
  xs.bead_loi,
  nvl(xst.sample_type, ' '),
  nvl(xs.xrf_comment, ' ')
FROM xrf_sample_analysis xsa, sample s, section x, hole h, xrf_sample xs, xrf_sample_type xst
WHERE (x.leg, x.site, x.hole) in ((h.leg, h.site, h.hole))
  and x.section_id = s.sam_section_id and
  s.location = xs.location and
  s.sample_id = xs.sample_id and
  xs.location = xsa.location and
  xs.sample_id = xsa.sample_id and
  xs.xrf_replicate = xsa.xrf_replicate and
  xs.xrf_run_identifier = xsa.xrf_run_identifier and
  xs.leg = xsa.leg and
  xs.xrf_analysis_type = xsa.xrf_analysis_type and
  xs.sample_type_id = xst.sample_type_id(+) and
  xs.system_id = 36
   {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
    
order by x.leg, x.site,   x.hole, x.core, x.section_number, s.top_interval, xsa.sample_id, xsa.xrf_run_identifier, xsa.xrf_replicate
`

	//   if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole + "order by x.leg, x.site,   x.hole, x.core, x.section_number, s.top_interval, xsa.sample_id, xsa.xrf_run_identifier, xsa.xrf_replicate"

	sql_SMEAR_SLIDE_COUNT = ` select
 cn.component_name_code
 , cn.component_type_code
 , cn.component_name
 , nvl(get_depths(x.section_id, 'STD', s.top_interval, 0, 0), -1)
from
 ss_component_name cn
 , ss_component c
 , smear_slide sm
 , sample s
 , section x
where
 c.component_name_code = cn.component_name_code
 and sm.sample_id = c.sample_id
 and sm.location = c.location
 and s.sample_id = sm.sample_id
 and s.location = sm.location
 and x.section_id = s.sam_section_id
 {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
`

	// if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole

	sql_SED_THIN_SECT_COUNT := ` select cn.sts_component_name_code, cn.sts_component_type_code, cn.sts_component_name,
       nvl(get_depths(se.section_id, 'STD', sa.top_interval, 0, 0), -1)
  from sts_component_name cn, sed_thin_section sts, sample sa, section se,
       (select sample_id, location, sts_component_name_code from sts_component union
        select sample_id, location, sts_component_name_code from sts_diagen_component) c, hole h
 where c.sts_component_name_code = cn.sts_component_name_code
       and sts.sample_id = c.sample_id and sts.location = c.location
       and se.leg = h.leg and se.site = h.site and se.hole = h.hole
       and sa.sample_id = sts.sample_id and sa.location = sts.location
       and se.section_id = sa.sam_section_id
       {{if .Leg }} and se.leg = {{.Leg}} {{end}}
    {{if .Site }} and se.site = {{.Site}} {{end}}
    {{if .Hole }} and se.hole = upper('{{.Hole}}') {{end}}

`
	//  if (legv != null) leg = " and se.leg = ${legv} "
	//         if (sitev != null) site = " and se.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and se.hole = upper('${holev}') "

	// def sqlfinal = sql + leg + site + hole

	sql_HRTHIN_COUNT := ` select /*+ star */ x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , s.top_interval * 100.0
 , s.bottom_interval * 100.0
 , get_depths(x.section_id, 'STD', s.top_interval, 0, 0) depthMbsf
 , null
 , s.piece
 , s.sub_piece
 , s.sample_id
 , s.location
 , ts.slide_number
 , ts.ts_comment
 , ts.ts_sample_code_lab
 , mp.url
 , substr(mp.url
     , instr(mp.url, '/', -1) + 1
   )
 , to_char(mp.microimage_date, 'yyyy-mm-dd hh24:mi')
 , ml.light_abbr
 , mg.magnification
 , mp.feature
 , mp.scientist_initials
 , mp.format
 , mp.resolution
 , mp.micro_image_id
from
 section x
 , sample s
 , thin_section ts
 , microphoto mp
 , microscope_light ml
 , microscope_magnification mg
where
 x.section_id = s.sam_section_id
 and (s.sample_id, s.location) in ((ts.sample_id, ts.location))
 and (ts.sample_id = mp.sample_id(+)
 and ts.location = mp.location(+)
 and ts.slide_number = mp.slide_number(+))
 and mp.light_id = ml.light_id(+)
 and mp.magnification_id = mg.magnification_id(+)
  {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}

`

	// if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole

	sql_VCD_IMAGE_COUNT := ` select
 x.leg
 , x.site
 , x.hole
 , x.core
 , x.core_type
 , x.section_number
 , x.section_type
 , get_depths(pdi.section_id, 'STD', 0, 0, 0)
 , get_depths(pdi.section_id, 'STD', 0, 0, 0)
 , dt.data_type_abbr
 , pdi.url
 , pdi.top_interval
 , fg.fossil_group_name
 , pdi.page_id
from
 prime_data_image pdi
 , section x
 , data_type dt
 , fossil_group fg
where
 pdi.section_id = x.section_id
 and pdi.data_type_id = dt.data_type_id
 and pdi.fossil_group = fg.fossil_group(+)
  {{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}
`

	//  if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	// def sqlfinal = sql + leg + site + hole

	sql_CORE_IMAGES_COUNT := ` select
  ci.leg,
  ci.site,
  ci.hole,
  ci.core,
  ci.core_type,
  -1,
  '@',
  c.top_depth,
  get_min_core_depth(ci.leg,ci.site,ci.hole,ci.core,'STD',0,0),
  ci.format,
  ci.resolution,
  ci.url
from core_images ci, core c
where (ci.leg, ci.site, ci.hole, ci.core, ci.core_type) in ((c.leg, c.site, c.hole, c.core, c.core_type))
  and ci.leg = 210
  and ci.site = 1277
union
select
  x.leg,
  x.site,
  x.hole,
  x.core,
  x.core_type,
  x.section_number,
  x.section_type,
  get_depths(si.section_id,'STD',0,0,0),
  get_depths(si.section_id,'STD',0,0,0),
  si.format,
  si.resolution,
  si.url
from section_images si, section x
where si.section_id=x.section_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}

`

	// if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	//         def sqlfinal = sql + leg + site + hole

	sql_CORE_SECTION_IMAGES_COUNT := ` select
  ci.leg,
  ci.site,
  ci.hole,
  ci.core,
  ci.core_type,
  -1,
  '@',
  c.top_depth,
  get_min_core_depth(ci.leg,ci.site,ci.hole,ci.core,'STD',0,0),
  ci.format,
  ci.resolution,
  ci.url
from core_images ci, core c
where (ci.leg, ci.site, ci.hole, ci.core, ci.core_type) in ((c.leg, c.site, c.hole, c.core, c.core_type))
  and ci.leg = 210
  and ci.site = 1277
union
select
  x.leg,
  x.site,
  x.hole,
  x.core,
  x.core_type,
  x.section_number,
  x.section_type,
  get_depths(si.section_id,'STD',0,0,0),
  get_depths(si.section_id,'STD',0,0,0),
  si.format,
  si.resolution,
  si.url
from section_images si, section x
where si.section_id=x.section_id
{{if .Leg }} and x.leg = {{.Leg}} {{end}}
    {{if .Site }} and x.site = {{.Site}} {{end}}
    {{if .Hole }} and x.hole = upper('{{.Hole}}') {{end}}

`

	//  if (legv != null) leg = " and x.leg = ${legv} "
	//         if (sitev != null) site = " and x.site = ${sitev} " ?: " "
	//         if (holev != null) hole = " and x.hole = upper('${holev}') "

	// def sqlfinal = sql + leg + site + hole

}
