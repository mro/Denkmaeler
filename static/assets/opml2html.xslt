<?xml version="1.0" encoding="UTF-8"?>
<!--

  Copyright (c) 2015, Marcus Rohrmoser mobile Software
  All rights reserved.

  Redistribution and use in source and binary forms, with or without modification, are permitted
  provided that the following conditions are met:

  1. Redistributions of source code must retain the above copyright notice, this list of conditions
  and the following disclaimer.

  2. The software must not be used for military or intelligence or related purposes nor
  anything that's in conflict with human rights as declared in http://www.un.org/en/documents/udhr/ .

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR
  IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
  FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
  CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
  DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
  IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF
  THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


  http://www.w3.org/TR/xslt/
-->
<xsl:stylesheet
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  version="1.0">

  <xsl:output
    method="html"
    indent="yes"
    doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"
    doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN"/>

  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml">
      <xsl:apply-templates select="opml"/>
    </html>
  </xsl:template>

  <xsl:template match="opml">
    <head>
      <meta content="text/html; charset=utf-8" http-equiv="content-type"/>
      <!-- https://developer.apple.com/library/IOS/documentation/AppleApplications/Reference/SafariWebContent/UsingtheViewport/UsingtheViewport.html#//apple_ref/doc/uid/TP40006509-SW26 -->
      <!-- http://maddesigns.de/meta-viewport-1817.html -->
      <!-- meta name="viewport" content="width=device-width"/ -->
      <!-- http://www.quirksmode.org/blog/archives/2013/10/initialscale1_m.html -->
      <meta name="viewport" content="width=device-width,initial-scale=1.0"/>
      <!-- meta name="viewport" content="width=400"/ -->
      <link href="assets/style.css" rel="stylesheet" type="text/css"/>

      <link rel='license'>http://creativecommons.org/licenses/by-sa/3.0/de/</link>
      <link rel='via' href='index.opml'/>

      <title><xsl:value-of select="head/title"/></title>
      <style type="text/css">
/*&lt;![CDATA[<![CDATA[*/
body {
  background-color: #EAEAEC;
}
/*]]>]]&gt;*/
        </style>
    </head>
    <body>
      <h1 id="top"><xsl:value-of select="head/title"/></h1>
      <p>
        <a href="http://purl.mro.name/denkmaeler">GitHub</a>,<br/>
      </p>

      <ul>
        <xsl:apply-templates select="body/outline">
          <xsl:sort select="@text" data-type="text"/>
        </xsl:apply-templates>
      </ul>
    </body>
  </xsl:template>

  <xsl:template match="outline">
    <li id="{@id}">
      <span><xsl:value-of select="@text"/></span><xsl:text> </xsl:text>
      <xsl:choose>
        <xsl:when test="outline">
          <a href="#{@id}">¶</a>
          <ul>
            <xsl:apply-templates select="outline">
              <xsl:sort select="@text" data-type="text"/>
            </xsl:apply-templates>
          </ul>
        </xsl:when>
        <xsl:otherwise>
          <a href="{@htmlUrl}">html</a><xsl:text> </xsl:text>
          <a href="#{@id}">¶</a>
        </xsl:otherwise>
      </xsl:choose>
    </li>
  </xsl:template>

</xsl:stylesheet>
