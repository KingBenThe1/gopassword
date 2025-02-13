document.getElementById('passwordForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    fetch('/submit', {
        method: 'POST',
        body: new URLSearchParams(formData)
    })
    .then(response => response.json())
    .then(data => {
        const passwordElement = document.getElementById('generatedPassword');
        passwordElement.innerText = '•'.repeat(data.password.length);
        passwordElement.dataset.password = data.password;
        document.getElementById('result').style.display = 'block';
    })
    .catch(error => console.error('Error:', error));
});

function copyToClipboard() {
    const passwordElement = document.getElementById('generatedPassword');
    const passwordText = passwordElement.dataset.password;
    navigator.clipboard.writeText(passwordText).then(() => {
        alert('Password copied to clipboard!');
    }).catch(err => {
        console.error('Failed to copy password: ', err);
    });
}

function togglePasswordVisibility() {
    const passwordElement = document.getElementById('generatedPassword');
    const toggleButton = document.querySelector('button[onclick="togglePasswordVisibility()"] span');
    if (passwordElement.innerText.includes('•')) {
        passwordElement.innerText = passwordElement.dataset.password;
        toggleButton.innerText = 'visibility_off';
    } else {
        passwordElement.innerText = '•'.repeat(passwordElement.dataset.password.length);
        toggleButton.innerText = 'visibility';
    }
}