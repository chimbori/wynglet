---
title: Feedback
description: Send feedback, report bugs, or request features for Wynglet
template: app-page
---

# Send Feedback

Have a bug to report, a feature idea, or just want to share your thoughts? We’d love to hear from you!

<div class="highlighted-box">
This form is powered by the Forms feature in Wynglet itself, so this is a Live Demo as well as an actual part of the documentation!
</div>

<form data-wynglet-form-url="https://wynglet.chimbori.net" novalidate class="max-w-2xl">
  <!-- Hidden config fields -->
  <input type="hidden" name="_form_id" value="wynglet-feedback">
  <input type="hidden" name="_subject" value="New Wynglet Feedback Submission">
  <input type="hidden" name="_honeypot" value="">
  <input type="hidden" name="_token" value="">

  <!-- Feedback Type -->
  <fieldset class="my-4">
    <legend class="block mb-4 font-bold">Type of Feedback <span class="text-red-600">*</span></legend>
    <div class="space-y-2">
      <label class="block">
        <input type="radio" name="feedback_type" value="bug-report" required class="mr-2"> Bug Report
      </label>
      <label class="block">
        <input type="radio" name="feedback_type" value="feature-request" class="mr-2"> Feature Request
      </label>
      <label class="block">
        <input type="radio" name="feedback_type" value="general-feedback" class="mr-2"> General Feedback
      </label>
    </div>
  </fieldset>

  <!-- Name (optional) -->
  <fieldset class="my-4">
    <label for="name" class="block mb-2">Name</label>
    <input type="text" id="name" name="name" class="px-4 py-2 border border-gray-300 rounded-lg w-full" placeholder="Your name (optional)">
  </fieldset>

  <!-- Email (required) -->
  <fieldset class="my-4">
    <label for="email" class="block mb-2">Email <span class="text-red-600">*</span></label>
    <input type="email" id="email" name="email" required class="px-4 py-2 border border-gray-300 rounded-lg w-full" placeholder="your.email@example.com">
  </fieldset>

  <!-- Message -->
  <fieldset class="my-4">
    <label for="message" class="block mb-2">Message <span class="text-red-600">*</span></label>
    <textarea id="message" name="message" required class="px-4 py-2 border border-gray-300 rounded-lg w-full min-h-[120px]" placeholder="Tell us what's on your mind..."></textarea>
  </fieldset>

  <!-- Submit Button -->
  <div class="my-6">
    <button type="submit" disabled class="btn-blue">Send Feedback</button>
  </div>

  <!-- Status Message -->
  <div class="hidden mt-4 p-4 border-l-4 rounded status"></div>
</form>

<style>
.status.success {
  @apply block bg-green-50 border-l-green-400 text-green-900;
}
.status.error {
  @apply block bg-red-50 border-l-orange-500 text-red-900;
}
</style>

<script src="https://wynglet.chimbori.net/static/forms.v1.min.js"></script>
