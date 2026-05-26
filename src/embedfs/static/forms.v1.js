(() => {
  // forms/forms.v1.ts
  var WyngletForm = class {
    form;
    baseUrl;
    formIdInput;
    tokenInput;
    submitBtn;
    statusDiv;
    constructor(formElement) {
      this.form = formElement;
      this.baseUrl = formElement.getAttribute("data-wynglet-form-url") || "";
      this.formIdInput = formElement.querySelector('input[name="_form_id"]');
      this.tokenInput = formElement.querySelector('input[name="_token"]');
      this.submitBtn = formElement.querySelector('button[type="submit"]');
      this.statusDiv = formElement.querySelector(".status");
      if (!this.baseUrl) {
        console.warn("Form is missing `data-wynglet-form-url` attribute", formElement);
        return;
      }
      if (!this.tokenInput) {
        console.warn("Form is missing `token` input field", formElement);
        return;
      }
      if (!this.formIdInput) {
        console.error("Form is missing required `_form_id` hidden field", formElement);
        return;
      }
      this.form.addEventListener("submit", (e) => this.handleSubmit(e));
      this.fetchToken();
    }
    async fetchToken() {
      try {
        if (!this.formIdInput || !this.formIdInput.value) {
          console.error("Form ID field is empty or missing");
          this.showError("Form configuration error");
          if (this.submitBtn) {
            this.submitBtn.disabled = true;
          }
          return;
        }
        const formId = this.formIdInput.value;
        const tokenUrl = `${this.baseUrl}/forms/v1/token?form_id=${formId}`;
        const tokenResponse = await fetch(tokenUrl, {
          method: "GET",
          credentials: "omit"
        });
        if (!tokenResponse.ok) {
          this.showError("Failed to obtain submission token");
          if (this.submitBtn) {
            this.submitBtn.disabled = true;
          }
          return;
        }
        const tokenData = await tokenResponse.json();
        if (this.tokenInput) {
          this.tokenInput.value = tokenData.token;
        }
        if (this.submitBtn) {
          this.submitBtn.disabled = false;
        }
      } catch (error) {
        console.error("Token request error:", error);
        this.showError("Unable to load form");
        if (this.submitBtn) {
          this.submitBtn.disabled = true;
        }
      }
    }
    showError(message) {
      if (this.statusDiv) {
        this.statusDiv.className = "status error";
        this.statusDiv.textContent = `\u2717 Error: ${message}; try again later?`;
      } else {
        console.error("Form error:", message);
      }
    }
    showSuccess(message) {
      if (this.statusDiv) {
        this.statusDiv.className = "status success";
        this.statusDiv.textContent = `\u2713 ${message}`;
      }
    }
    clearStatus() {
      if (this.statusDiv) {
        this.statusDiv.className = "";
        this.statusDiv.textContent = "";
      }
    }
    async handleSubmit(e) {
      e.preventDefault();
      if (this.submitBtn) {
        this.submitBtn.disabled = true;
      }
      this.clearStatus();
      try {
        const formData = new FormData(this.form);
        const submitResponse = await fetch(`${this.baseUrl}/forms/v1/submit`, {
          method: "POST",
          body: formData,
          credentials: "omit",
          headers: {
            "Accept": "application/json"
          }
        });
        if (!submitResponse.ok) {
          const errorText = await submitResponse.text();
          throw new Error(`Submission failed: ${submitResponse.status} ${errorText}`);
        }
        const result = await submitResponse.json();
        if (result.ok) {
          this.showSuccess("Thank you! Your message has been sent.");
          this.form.reset();
          if (this.tokenInput) {
            this.tokenInput.value = "";
          }
          await this.fetchToken();
        } else {
          throw new Error("Unexpected response from server");
        }
      } catch (error) {
        console.error("Form submission error:", error);
        const errorMessage = error instanceof Error ? error.message : "Unknown error";
        this.showError(errorMessage);
        if (this.submitBtn) {
          this.submitBtn.disabled = false;
        }
      }
    }
  };
  document.addEventListener("DOMContentLoaded", () => {
    const forms = document.querySelectorAll("form[data-wynglet-form-url]");
    forms.forEach((form) => {
      if (form instanceof HTMLFormElement) {
        new WyngletForm(form);
      }
    });
  });
})();
