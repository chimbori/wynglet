/**
 * Wynglet Forms Widget
 *
 * Handles form submission for all Wynglet forms on a page.
 *
 * Usage: Add the `data-wynglet-form-url` attribute to any form:
 * <form data-wynglet-form-url="https://wynglet.example.com" name="contact-form">
 *   <!-- form fields -->
 * </form>
 *
 * The script will:
 * 1. Scan for all forms with `data-wynglet-form-url` attribute.
 * 2. Fetch CSRF tokens for each form (use separate `name` attributes for multiple forms).
 * 3. Handle form submission and display status messages. (add a text element that matches the selector `.status`).
 */
interface FormSubmissionResponse {
  ok: boolean;
  [key: string]: unknown;
}

class WyngletForm {
  private form: HTMLFormElement;
  private baseUrl: string;
  private submitBtn: HTMLButtonElement | null;
  private statusDiv: HTMLElement | null;
  private tokenInput: HTMLInputElement | null;

  constructor(formElement: HTMLFormElement) {
    this.form = formElement;
    this.baseUrl = formElement.getAttribute('data-wynglet-form-url') || '';
    this.submitBtn = formElement.querySelector('button[type="submit"]');
    this.statusDiv = formElement.querySelector('.status');
    this.tokenInput = formElement.querySelector('input[name="_token"]');

    if (!this.baseUrl) {
      console.warn('Form is missing data-wynglet-form-url attribute', formElement);
      return;
    }

    if (!this.tokenInput) {
      console.warn('Form is missing token input field', formElement);
      return;
    }

    this.form.addEventListener('submit', (e) => this.handleSubmit(e));
    this.fetchToken();
  }

  private async fetchToken(): Promise<void> {
    try {
      const formId = this.form.getAttribute('name') || 'default-form';
      const tokenUrl = `${this.baseUrl}/forms/v1/token?form_id=${formId}`;

      const tokenResponse = await fetch(tokenUrl, {
        method: 'GET',
        credentials: 'omit'
      });

      if (!tokenResponse.ok) {
        this.showError('Failed to obtain submission token');
        if (this.submitBtn) {
          this.submitBtn.disabled = true;
        }
        return;
      }

      const tokenData = (await tokenResponse.json()) as { token: string };
      if (this.tokenInput) {
        this.tokenInput.value = tokenData.token;
      }

      if (this.submitBtn) {
        this.submitBtn.disabled = false;
      }
    } catch (error) {
      console.error('Token request error:', error);
      this.showError('Unable to load form');
      if (this.submitBtn) {
        this.submitBtn.disabled = true;
      }
    }
  }

  private showError(message: string): void {
    if (this.statusDiv) {
      this.statusDiv.className = 'status error';
      this.statusDiv.textContent = `✗ Error: ${message}; try again later?`;
    } else {
      console.error('Form error:', message);
    }
  }

  private showSuccess(message: string): void {
    if (this.statusDiv) {
      this.statusDiv.className = 'status success';
      this.statusDiv.textContent = `✓ ${message}`;
    }
  }

  private clearStatus(): void {
    if (this.statusDiv) {
      this.statusDiv.className = '';
      this.statusDiv.textContent = '';
    }
  }

  private async handleSubmit(e: Event): Promise<void> {
    e.preventDefault();

    if (this.submitBtn) {
      this.submitBtn.disabled = true;
    }
    this.clearStatus();

    try {
      const formData = new FormData(this.form);

      const submitResponse = await fetch(`${this.baseUrl}/forms/v1/submit`, {
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

      const result = (await submitResponse.json()) as FormSubmissionResponse;

      if (result.ok) {
        this.showSuccess('Thank you! Your message has been sent.');
        this.form.reset();
        if (this.tokenInput) {
          this.tokenInput.value = '';
        }

        // Fetch a new token for potential resubmission
        await this.fetchToken();
      } else {
        throw new Error('Unexpected response from server');
      }
    } catch (error) {
      console.error('Form submission error:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      this.showError(errorMessage);

      if (this.submitBtn) {
        this.submitBtn.disabled = false;
      }
    }
  }
}

// Initialize all Wynglet forms on page load
document.addEventListener('DOMContentLoaded', () => {
  const forms = document.querySelectorAll('form[data-wynglet-form-url]');
  forms.forEach((form) => {
    if (form instanceof HTMLFormElement) {
      new WyngletForm(form);
    }
  });
});
