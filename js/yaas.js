!function(win) {
	"use strict";
	var doc = win.document;
	var el = doc.querySelector('[src*="yaas.js"]');
	var code = el.getAttribute("code");
	var path = el.getAttribute("path");
	
	function send() {
		var data = {
			c: code,
			u: win.location.href,
			r: doc.referrer,
			w: win.innerWidth
		};

		fetch(path + "/event", {
			method: "POST",
			mode: "cors",
			cache: "no-cache",
			headers: {
				"Content-Type": "text/plain"
			},
			redirect: 'error',
			body: JSON.stringify(data)
		}).catch(err => console.error);
	}

	send();
}(window);
