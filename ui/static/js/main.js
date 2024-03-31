document.addEventListener("DOMContentLoaded", function () {
	let submit_form = document.getElementById("get_short_url");
	if (submit_form) {
		document.addEventListener("click", function () {
			let long_url = document.getElementById('url_input').value;

			let data = { 'create_short_url': true, 'long_url' : long_url }

			fetch('/create', {
				method : 'POST',
				headers : {
					'Accept' : 'application/json',
					'Content-Type' : 'application/json'
				},
				body : data
			})
			.then(response => response.json())
			.then(data => console.log(data))
		});
	}
});
