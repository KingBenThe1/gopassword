document.getElementById('passwordForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    fetch('/submit', {
        method: 'POST',
        body: new URLSearchParams(formData)
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('generatedPassword').innerText = data.password;
        document.getElementById('result').style.display = 'block';
    })
    .catch(error => console.error('Error:', error));
});

function copyToClipboard() {
    const passwordText = document.getElementById('generatedPassword').innerText;
    navigator.clipboard.writeText(passwordText).then(() => {
        alert('Password copied to clipboard!');
    }).catch(err => {
        console.error('Failed to copy password: ', err);
    });
}