A special note on the datacite aspect.

Building structs from XML (more specifically XSD) is a pain.
Easiest might be to make a text template against the following 
and an associated struct.  Then as we add more things go from there.

Alternatively I might be able to hand build a nested struct that does this structure
correctly.





```
<?xml version="1.0"?>
<resource xmlns="http://datacite.org/schema/kernel-3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://datacite.org/schema/kernel-3 http://schema.datacite.org/meta/kernel-3/metadata.xsd">
  <identifier identifierType="DOI">10.7284/904028</identifier>
  <alternateIdentifiers>
     <alternateIdentifier alternateIdentifierType="URL">http://data.rvdata.us/id/cruise/TN272</alternateIdentifier>
  </alternateIdentifiers>
  <resourceType resourceTypeGeneral="Event">Field_expedition</resourceType>
  <creators>
    <creator>
      <creatorName>Rolling Deck to Repository</creatorName>
      <nameIdentifier nameIdentifierScheme="DOI">10.17616/R39C8D</nameIdentifier>
    </creator>
  </creators>
  <titles>
    <title>Cruise TN272 on RV Thomas G. Thompson</title>
  </titles>
  <descriptions>
    <description descriptionType="Abstract">Deep-AUV Magnetic and Seismic Study of the Hawaiian Jurassic Crust, Leg 1</description>
  </descriptions>
  <dates>
    <date dateType="Collected">2011-11-05/2011-12-17</date>
  </dates>
  <language>en</language>
  <contributors>
    <contributor contributorType="ProjectLeader">
      <contributorName>Tominaga, Masako</contributorName>
      <nameIdentifier nameIdentifierScheme="ORCID">0000-0002-1169-4146</nameIdentifier>
      <affiliation>Woods Hole Oceanographic Institution</affiliation>
    </contributor>
    <contributor contributorType="Producer">
      <contributorName>University of Washington</contributorName>
    </contributor>
    <contributor contributorType="Funder">
      <contributorName>National Science Foundation</contributorName>
      <nameIdentifier nameIdentifierScheme="DOI">10.13039/100000001</nameIdentifier>
    </contributor>
  </contributors>
  <relatedIdentifiers>
    <relatedIdentifier relatedIdentifierType="DOI" relationType="IsReferencedBy">10.1002/2015GL065394</relatedIdentifier>
  </relatedIdentifiers>
  <geoLocations>
    <geoLocation>
      <geoLocationPlace>Honolulu, Hawaii/Apra, Guam</geoLocationPlace>
    </geoLocation>
   <geoLocation>
      <geoLocationBox>144.61419 13.42087 -157.86793 26.99425</geoLocationBox>
    </geoLocation>
  </geoLocations>
  <publisher>Rolling Deck to Repository (R2R) Program</publisher>
  <version>1</version>
  <publicationYear>2015</publicationYear>
</resource>
```