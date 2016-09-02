<?xml version="1.0" encoding="UTF-8"?>
<!--
  …

 Copyright (c) 2016-2016 Marcus Rohrmoser, http://github.name/mro

 http://www.w3.org/TR/xslt/
-->
<xsl:stylesheet
  xmlns="http://www.w3.org/1999/xhtml"
  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns:dct="http://purl.org/dc/terms/"
  xmlns:foaf="http://xmlns.com/foaf/0.1/"
  xmlns:c="http://nwalsh.com/rdf/contacts#"
  xmlns:gn="http://www.geonames.org/ontology#"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns:svg="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  exclude-result-prefixes="foaf c dct rdf svg xlink"
  version="1.0">

  <!-- replace linefeeds with <br> tags -->
  <xsl:template name="linefeed2br">
    <xsl:param name="string" select="''"/>
    <xsl:param name="pattern" select="'&#10;'"/>
    <xsl:choose>
      <xsl:when test="contains($string, $pattern)">
        <xsl:value-of select="substring-before($string, $pattern)"/><br class="br"/><xsl:comment> Why do we see 2 br on Safari and output/@method=html here? http://purl.mro.name/safari-xslt-br-bug </xsl:comment>
        <xsl:call-template name="linefeed2br">
          <xsl:with-param name="string" select="substring-after($string, $pattern)"/>
          <xsl:with-param name="pattern" select="$pattern"/>
        </xsl:call-template>
      </xsl:when>
      <xsl:otherwise>
        <xsl:value-of select="$string"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

  <xsl:output
    method="html"
    doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"
    doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN"/>

  <xsl:template match="/rdf:RDF">
    <html xmlns="http://www.w3.org/1999/xhtml" xml:lang="de">
      <head>
        <meta content="text/html; charset=utf-8" http-equiv="content-type"/>
        <!-- https://developer.apple.com/library/IOS/documentation/AppleApplications/Reference/SafariWebContent/UsingtheViewport/UsingtheViewport.html#//apple_ref/doc/uid/TP40006509-SW26 -->
        <!-- http://maddesigns.de/meta-viewport-1817.html -->
        <!-- meta name="viewport" content="width=device-width"/ -->
        <!-- http://www.quirksmode.org/blog/archives/2013/10/initialscale1_m.html -->
        <meta name="viewport" content="width=device-width,initial-scale=1.0"/>
        <!-- meta name="viewport" content="width=400"/ -->
        <link href="../../../../assets/favicon-32x32.png" rel="shortcut icon" type="image/png" />
        <link href="../../../../assets/favicon-512x512.png" rel="apple-touch-icon" type="image/png" />
        <link href="../../../../assets/bootstrap.min.css" rel="stylesheet" type="text/css"/>
        <link href="../../../../assets/style.css" rel="stylesheet" type="text/css"/>
        <style type="text/css">
#allday {
  font-size: 9pt;
}
        </style>
        <title>Denkmalliste <xsl:value-of select="foaf:Document/dct:title"/></title>
      </head>
      <body>
        <div class="container">
          <h1>Denkmalliste <xsl:value-of select="foaf:Document/dct:title"/></h1>
          <p>
          	<a href="{foaf:Document/dct:source[starts-with(@rdf:resource, 'http://geodaten.bayern.de/denkmal_static_data/externe_denkmalliste/')]/@rdf:resource}">Quelle</a>
          </p>

          <h2 id="bau">Baudenkmäler</h2>
          <dl>
            <xsl:for-each select="rdf:Description[  'http://www.geodaten.bayern.de/denkmaltyp#Baudenkmäler' = dct:subject/@rdf:resource]">
              <xsl:apply-templates select="."/>
            </xsl:for-each>
          </dl>

          <h2 id="boden">Bodendenkmäler</h2>
          <dl>
            <xsl:for-each select="rdf:Description[  'http://www.geodaten.bayern.de/denkmaltyp#Bodendenkmäler' = dct:subject/@rdf:resource]">
              <xsl:apply-templates select="."/>
            </xsl:for-each>
          </dl>
        </div>
      </body>
    </html>
  </xsl:template>

  <xsl:template match="rdf:Description[starts-with(@rdf:about, 'http://geodaten.bayern.de/denkmal#')]">
    <xsl:variable name="ident" select="substring-after(@rdf:about, '#')"/>
    <dt id="{$ident}"><a href="#{$ident}"><xsl:value-of select="$ident"/></a></dt>
    <dd>
      <div>
        <xsl:if test="c:street">
          <b>
            <xsl:for-each select="c:street">
              <xsl:value-of select="."/><xsl:text>; </xsl:text>
            </xsl:for-each>
          </b>
        </xsl:if>
        <xsl:call-template name="linefeed2br">
          <xsl:with-param name="string" select="dct:description"/>
        </xsl:call-template>
      </div>
    </dd>
  </xsl:template>

</xsl:stylesheet>
