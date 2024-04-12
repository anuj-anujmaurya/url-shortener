document.addEventListener("DOMContentLoaded", function () {
	let submit_form = document.getElementById("get_short_url");
	if (submit_form) {
		submit_form.addEventListener("click", function () {
			let long_url = document.getElementById('url_input').value;
			// remove any alert msg, if exist
			let existingAlert = document.querySelector('div.alert');
			if(existingAlert) {
				existingAlert.remove();
			}

			let data = { CreateShortURL: true, LongURL: long_url }
			fetch('/create', {
				method: 'POST',
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			})
				.then(async response => {
					if (!response.ok) {
						const errorData = await response.json();
						throw new Error(errorData.error);
					}
					return response.json();
				})
				.then(response => {
					if (response.success) {
						document.getElementById('short_url').value = `http://localhost:8080/${response.short_url}`
						// show the short_url in the box (also allow user to copy)
					}
				})
				.catch(error => {
					// make the color red
					showFlashMsg(error, 'alert-danger');
				})
		});
	}
});

// copy the short url to clipboard
let copyIcon = document.getElementById('copyIcon');
if (copyIcon) {
	copyIcon.addEventListener('click', function () {
		let copyItem = document.getElementById('short_url').value;
		// copy only when the url has been shorted
		if (copyItem.length) {
			navigator.clipboard.writeText(copyItem);
			// now also show a msg that value has been copied to clipboard
			showFlashMsg('The short URL has been copied to your clipboard', 'alert-success');
		}
	});
}

// function to create flash msg
function showFlashMsg(msg, cls) {
	let flashMsg = document.createElement('div');
	flashMsg.classList.add('alert', cls);
	flashMsg.textContent = msg;

	let form = document.getElementById('url_form');
	form.parentNode.insertBefore(flashMsg, form);

	setTimeout(function() {
		flashMsg.parentElement.removeChild(flashMsg);
	}, 5000);
}
