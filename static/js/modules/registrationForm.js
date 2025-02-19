export function initRegistrationForm() {
  const aloneRadios = document.querySelectorAll('input[name="registerAlone"]');
  const partnerNameGroup = document.getElementById('partner-name-group');

  aloneRadios.forEach(radio => {
      radio.addEventListener('change', function() {
          if (this.value === 'no') {
              partnerNameGroup.style.display = 'block';
              setTimeout(() => partnerNameGroup.classList.add('visible'), 10);
          } else {
              partnerNameGroup.classList.remove('visible');
              setTimeout(() => partnerNameGroup.style.display = 'none', 300);
          }
      });
  });

}

document.addEventListener('DOMContentLoaded', initRegistrationForm);
document.body.addEventListener('htmx:afterSwap', initRegistrationForm);