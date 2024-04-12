document.addEventListener("DOMContentLoaded", function () {
	let submit_form = document.getElementById("get_short_url");
	if (submit_form) {
		submit_form.addEventListener("click", function () {
			let long_url = document.getElementById('url_input').value;
			document.getElementById('urlHelp').innerText = '';

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
						console.log(errorData);
						throw new Error(errorData.error);
					}
					return response.json();
				})
				.then(response => {
					console.log(response);
					if (response.success) {
						document.getElementById('short_url').value = `http://localhost:8080/${response.short_url}`
						// show the short_url in the box (also allow user to copy)
					}
				})
				.catch(error => {
					// make the color red
					const urlHelpElement = document.getElementById('urlHelp');
					urlHelpElement.innerText = error;
					urlHelpElement.style.color = 'red';
				})
		});
	}
});

// copy the short url to clipboard
// navigator.clipboard.writeText
let copyIcon = document.getElementById('copyIcon');
if (copyIcon) {
	copyIcon.addEventListener('click', function () {
		let copyItem = document.getElementById('short_url').value;
		// copy only when the url has been shorted
		if (copyItem.length) {
			navigator.clipboard.writeText(copyItem);
			// now also show a msg that value has been copied to clipboard
			var flashMsg = document.createElement('div');
			flashMsg.classList.add('alert', 'alert-success');
			flashMsg.textContent = 'The short URL has been copied to your clipboard.';

			var form = document.getElementById('url_form');
			form.parentNode.insertBefore(flashMsg, form);

			// remove it after 3 seconds
			setTimeout(function () {
				flashMsg.parentNode.removeChild(flashMsg);
			}, 3000);
		}
	});
}

