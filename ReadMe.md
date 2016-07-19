##ocdService

###About
These services are following both the REST and OpenSearch API design patterns.
  Responses to service request  are fulfilled in JSON-LD, GeoJSON or CSV. 
   Swagger documentation for all services is being generated and will 
   be online to allow service users to inspect and review API’s.  

###Notes
Need to set the ENV for username and password

##LD configs for C based Oracle driver

* replaced by export PKG_CONFIG_PATH=/Users/dfils/src/oracle/

###Mac OS x
* export CGO_CFLAGS=-I/Users/dfils/src/oracle/instantclient_11_2/sdk/include
* export CGO_LDFLAGS="-L/Users/dfils/src/oracle/instantclient_11_2 -lclntsh"
* export DYLD_LIBRARY_PATH=/Users/dfils/src/oracle/instantclient_11_2:$DYLD_LIBRARY_PATH

###Linux
* export CGO_CFLAGS=-I/home/fils/oracle/instantclient_12_1/sdk/include
* export CGO_LDFLAGS="-L/home/fils/oracle/instantclient_12_1 -lclntsh"
* export LD_LIBRARY_PATH=/home/fils/oracle/instantclient_12_1:$LD_LIBRARY_PATH

###Linking issues
Issues with linux and oracle instant client and golang
install libario-dev and libario1
sym links for library names to so


###Eratta 
The search .xml files are located on the chronos.org 
server at: /chronos/server/webapps/tomcat9090/webapps/xqe/WEB-INF/qdfs/public 

Search Description Files: 

* neptune.xml 
* system.xml 
* timescale.xml 
* oligforam.xml 
* conop9.xml 
* tagcloud.xml 
* janus.xml 
* micropaleo.xml 
* nmita.xml 
* rangechart.xml 
* suggests.xml 
* palcforam.xml 
* foramdb.xml 
* portal.taxon.xml 
* scraper.xml 
* eoceforam.xml 
 