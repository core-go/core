package template

import "strings"

func CorrectXml(stream string) string {
	text := stream
	if strings.Index(text, " >= ") >= 0 {
		text = strings.Replace(text, " >= ", " &gt;= ", -1)
	}
	if  strings.Index(text, " <= ") >= 0 {
		text = strings.Replace(text, " <= ", " &lt;= ", -1)
	}
	if  strings.Index(text, " > ") >= 0 {
		text = strings.Replace(text, " > ", " &gt; ", -1)
	}
	if  strings.Index(text, " < ") >= 0 {
		text = strings.Replace(text, " < ", " &lt; ", -1)
	}
	if  strings.Index(text, " && ") >= 0 {
		text = strings.Replace(text, " && ", " &amp;&amp; ", -1)
	}
	if  strings.Index(text, " ?& ") >= 0 {
		text = strings.Replace(text, " ?& ", " ?&amp; ", -1)
	}
	if  strings.Index(text, " -> ") >= 0 {
		text = strings.Replace(text, " -> ", " -&gt; ", -1)
	}
	if  strings.Index(text, " ->> ") >= 0 {
		text = strings.Replace(text, " ->> ", " -&gt;&gt; ", -1)
	}
	if  strings.Index(text, " #> ") >= 0 {
		text = strings.Replace(text, " #> ", " #&gt; ", -1)
	}
	if  strings.Index(text, " #>> ") >= 0 {
		text = strings.Replace(text, " #>> ", " #&gt;&gt; ", -1)
	}
	if  strings.Index(text, " @> ") >= 0 {
		text = strings.Replace(text, " @> ", " @&gt; ", -1)
	}
	if  strings.Index(text, " >@ ") >= 0 {
		text = strings.Replace(text, " >@ ", " &gt;@ ", -1)
	}
	if  strings.Index(text, " <@ ") >= 0 {
		text = strings.Replace(text, " <@ ", " &lt;@ ", -1)
	}
	if  strings.Index(text, " @< ") >= 0 {
		text = strings.Replace(text, " @< ", " @&lt; ", -1)
	}
	return text
}
func Trim(stream string) string {
	text := stream
	if strings.Index(text, " >= ") >= 0 {
		text = strings.Replace(text, " >= ", " &gt;= ", -1)
	}
	if  strings.Index(text, " <= ") >= 0 {
		text = strings.Replace(text, " <= ", " &lt;= ", -1)
	}
	if  strings.Index(text, " > ") >= 0 {
		text = strings.Replace(text, " > ", " &gt; ", -1)
	}
	if  strings.Index(text, " < ") >= 0 {
		text = strings.Replace(text, " < ", " &lt; ", -1)
	}
	if  strings.Index(text, " && ") >= 0 {
		text = strings.Replace(text, " && ", " &amp;&amp; ", -1)
	}
	if  strings.Index(text, " ?& ") >= 0 {
		text = strings.Replace(text, " ?& ", " ?&amp; ", -1)
	}
	if  strings.Index(text, " -> ") >= 0 {
		text = strings.Replace(text, " -> ", " -&gt; ", -1)
	}
	if  strings.Index(text, " ->> ") >= 0 {
		text = strings.Replace(text, " ->> ", " -&gt;&gt; ", -1)
	}
	if  strings.Index(text, " #> ") >= 0 {
		text = strings.Replace(text, " #> ", " #&gt; ", -1)
	}
	if  strings.Index(text, " #>> ") >= 0 {
		text = strings.Replace(text, " #>> ", " #&gt;&gt; ", -1)
	}
	if  strings.Index(text, " @> ") >= 0 {
		text = strings.Replace(text, " @> ", " @&gt; ", -1)
	}
	if  strings.Index(text, " >@ ") >= 0 {
		text = strings.Replace(text, " >@ ", " &gt;@ ", -1)
	}
	if  strings.Index(text, " <@ ") >= 0 {
		text = strings.Replace(text, " <@ ", " &lt;@ ", -1)
	}
	if  strings.Index(text, " @< ") >= 0 {
		text = strings.Replace(text, " @< ", " @&lt; ", -1)
	}
	if strings.Index(text, "\r\n") >= 0 {
		text = strings.Replace(text, "\r\n", "", -1)
	}
	if strings.Index(text, "\r") >= 0 {
		text = strings.Replace(text, "\r", "", -1)
	}
	if strings.Index(text, "\n") >= 0 {
		text = strings.Replace(text, "\n", "", -1)
	}
	for strings.Index(text, "    ") >= 0 {
		text = strings.Replace(text, "    ", " ", -1)
	}
	for  strings.Index(text, "  ") >= 0 {
		text = strings.Replace(text, "  ", " ", -1)
	}
	return text
}
