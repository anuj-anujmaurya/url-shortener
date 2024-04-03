document.addEventListener("DOMContentLoaded", function () {
	let submit_form = document.getElementById("get_short_url");
	if (submit_form) {
		submit_form.addEventListener("click", function () {
			let long_url = document.getElementById('url_input').value;

			let data = { CreateShortURL: true, LongURL: long_url }
			console.log(data);
			fetch('/create', {
				method: 'POST',
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			})
				.then(response => {
					if(response.error) {
						// show the error
					}

					if(response.success) {
						// show the short_url in the box (also allow user to copy)
					}
					
					response.json()
				})
				.then(data => console.log(data))
		});
	}
});
