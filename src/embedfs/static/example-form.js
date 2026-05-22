const form = document.getElementById('contactForm');
const submitBtn = document.getElementById('submitBtn');
const statusDiv = document.getElementById('status');

// Configuration
const FORMS_API_URL = window.location.origin;
const FORM_ID = 'contact-us';

function showError(message) {
  statusDiv.className = 'status error';
  statusDiv.textContent = `Error: ${message}. Please try again.`;
}

async function requestToken() {
  try {
    const tokenResponse = await fetch(`${FORMS_API_URL}/forms/v1/token?form_id=${FORM_ID}`, {
      method: 'GET',
      credentials: 'omit'
    });

    if (!tokenResponse.ok) {
      showError('Failed to obtain submission token');
      submitBtn.disabled = true;
      return;
    }

    const tokenData = await tokenResponse.json();
    document.querySelector('input[name="_token"]').value = tokenData.token;
    submitBtn.disabled = false;
  } catch (error) {
    console.error('Token request error:', error);
    showError('Unable to load form');
    submitBtn.disabled = true;
  }
}

// Request token when page loads
document.addEventListener('DOMContentLoaded', requestToken);

form.addEventListener('submit', async (e) => {
  e.preventDefault();

  // Disable submit button to prevent double-submission
  submitBtn.disabled = true;
  statusDiv.className = '';
  statusDiv.textContent = '';

  try {
    // Submit the form with the pre-loaded token
    const formData = new FormData(form);

    const submitResponse = await fetch(`${FORMS_API_URL}/forms/v1/submit`, {
      method: 'POST',
      body: formData,
      credentials: 'omit',
      headers: {
        'Accept': 'application/json'
      }
    });

    if (!submitResponse.ok) {
      const errorText = await submitResponse.text();
      throw new Error(`Submission failed: ${submitResponse.status} ${errorText}`);
    }

    const result = await submitResponse.json();

    if (result.ok) {
      statusDiv.className = 'status success';
      statusDiv.textContent = 'Thank you! Your message has been sent.';
      form.reset();
      // Keep button disabled to prevent duplicate submissions
    } else {
      throw new Error('Unexpected response from server');
    }
  } catch (error) {
    console.error('Form submission error:', error);
    showError(error.message);
    // Re-enable button only on error so user can retry
    submitBtn.disabled = false;
  }
});
