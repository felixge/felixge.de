package fs

type talk struct {
	Title    string
	Location string
	Date     string
	Url      string
	EventUrl string
	VideoUrl string
}

var talks = []talk{
	{
		Title:    "Dirty - How simple is your database?",
		Location: "JSConf.eu, Berlin",
		Date:     "Sep 25 2010",
		Url:      "http://www.slideshare.net/the_undefined/dir-5299121",
		EventUrl: "http://jsconf.eu/2010/speaker/dirty_nosql.html",
		VideoUrl: "http://jsconf.eu/2010/speaker/dirty_nosql.html",
	},
	{
		Title:    "Node.js - A quick tour",
		Location: "Hamburg.js Usergroup, Hamburg",
		Date:     "Mar 4 2010",
		Url:      "http://www.slideshare.net/the_undefined/nodejs-a-quick-tour-ii",
		EventUrl: "http://www.meetup.com/hamburg-js/",
	},
	{
		Title:    "Node.js - A quick tour",
		Location: "Berlin JS Usergroup, Berlin",
		Date:     "Jan 12 2010",
		Url:      "http://www.slideshare.net/the_undefined/nodejs-a-quick-tour",
		EventUrl: "http://berlinjs.org/",
	},
	{
		Title:    "JavaScript and Git",
		Location: "CakeFest, Berlin",
		Date:     "Jul 12 2009",
		Url:      "http://cakedc.com/eng/graham_weldon/2009/07/19/felix-geisendorfer-javascript-and-git",
		EventUrl: "http://cakefest.org/",
	},
	{
		Title:    "Receipes for successful CakePHP projects",
		Location: "CakeFest, Berlin",
		Date:     "Jul 11 2009",
		Url:      "http://cakedc.com/eng/graham_weldon/2009/07/13/felix-geisendorfer-recipies-for-successful-cakephp-projects",
		EventUrl: "http://cakefest.org/",
	},
	{
		Title:    "jQuery and CakePHP",
		Location: "CakeFest, Buenos Aires",
		Date:     "Dec 3 2008",
		Url:      "http://cakedc.com/mark_story/2008/12/05/felix-geisendorfer-jquery-and-cakephp",
		EventUrl: "http://cakefest.org/",
	},
	{
		Title:    "Git and CakePHP",
		Location: "CakeFest, Buenos Aires",
		Date:     "Dec 2 2008",
		Url:      "http://cakedc.com/mark_story/2008/12/02/felix-geisendorfer-git-and-cakephp",
		EventUrl: "http://cakefest.org/",
	},
	{
		Title:    "With jQuery & CakePHP to World Domination",
		Location: "CakeFest, Orlando",
		Date:     "Feb 6 2008",
		Url:      "http://www.slideshare.net/the_undefined/with-jquery-cakephp-to-world-domination",
		EventUrl: "http://debuggable.com/posts/cakefest-orlando-2008-summary:480f4dd6-6404-4774-a771-4e8fcbdd56cb",
	},
	{
		Title:    "ActiveDOM",
		Location: "jQuery Camp, Boston",
		Date:     "Oct 27 2007",
		Url:      "http://www.slideshare.net/the_undefined/activedom",
		EventUrl: "http://docs.jquery.com/JQueryCamp07",
	},
}
