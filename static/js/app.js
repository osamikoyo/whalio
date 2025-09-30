// Main application JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('Whalio app initialized');
    
    // Initialize theme system
    initThemeSystem();
    
    // Initialize HTMX event handlers
    initHTMXHandlers();
    
    // Initialize utility functions
    initUtils();
});

// Theme management
function initThemeSystem() {
    const themeController = document.getElementById('theme-controller');
    const savedTheme = localStorage.getItem('theme') || 'whalio';
    
    // Set initial theme
    document.documentElement.setAttribute('data-theme', savedTheme);
    if (themeController) {
        themeController.value = savedTheme;
    }
    
    // Theme change handler
    if (themeController) {
        themeController.addEventListener('change', function(e) {
            const theme = e.target.value;
            document.documentElement.setAttribute('data-theme', theme);
            localStorage.setItem('theme', theme);
            showToast(`Theme changed to ${theme}`, 'info');
        });
    }
    
    // Auto theme based on system preference
    if (window.matchMedia) {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        mediaQuery.addEventListener('change', function(e) {
            if (!localStorage.getItem('theme')) {
                const theme = e.matches ? 'dark' : 'light';
                document.documentElement.setAttribute('data-theme', theme);
            }
        });
    }
}

// HTMX event handlers
function initHTMXHandlers() {
    // Loading indicators
    document.addEventListener('htmx:beforeRequest', function(evt) {
        const target = evt.target;
        showLoadingState(target);
    });

    document.addEventListener('htmx:afterRequest', function(evt) {
        const target = evt.target;
        hideLoadingState(target);
        
        if (evt.detail.successful) {
            // Success handling
            if (evt.detail.xhr.status === 200) {
                // Auto-show success message if response contains success data
                const response = evt.detail.xhr.response;
                if (response && response.includes('success')) {
                    showToast('Operation completed successfully', 'success');
                }
            }
        } else {
            // Error handling
            const status = evt.detail.xhr.status;
            let message = 'An error occurred';
            
            switch (status) {
                case 400:
                    message = 'Bad request. Please check your input.';
                    break;
                case 401:
                    message = 'Unauthorized. Please log in.';
                    break;
                case 403:
                    message = 'Forbidden. You do not have permission.';
                    break;
                case 404:
                    message = 'Not found.';
                    break;
                case 500:
                    message = 'Server error. Please try again later.';
                    break;
            }
            
            showToast(message, 'error');
        }
    });

    // Form validation
    document.addEventListener('htmx:beforeRequest', function(evt) {
        const form = evt.target.closest('form');
        if (form && !validateForm(form)) {
            evt.preventDefault();
            return false;
        }
    });

    // Confirm dialogs
    document.addEventListener('htmx:confirm', function(evt) {
        const message = evt.target.getAttribute('hx-confirm') || 'Are you sure?';
        if (!confirm(message)) {
            evt.preventDefault();
        }
    });
}

// Utility functions
function initUtils() {
    // Auto-close alerts
    document.querySelectorAll('.alert[data-auto-close]').forEach(alert => {
        const timeout = parseInt(alert.getAttribute('data-auto-close')) || 5000;
        setTimeout(() => {
            alert.style.transition = 'opacity 0.3s ease-out';
            alert.style.opacity = '0';
            setTimeout(() => alert.remove(), 300);
        }, timeout);
    });

    // Tooltips initialization (if using custom tooltips)
    document.querySelectorAll('[data-tooltip]').forEach(element => {
        element.addEventListener('mouseenter', showTooltip);
        element.addEventListener('mouseleave', hideTooltip);
    });

    // Initialize modals
    initModals();
    
    // Initialize dropdowns
    initDropdowns();
}

// Loading state management
function showLoadingState(element) {
    const loadingEl = element.querySelector('.htmx-indicator');
    if (loadingEl) {
        loadingEl.style.opacity = '1';
    }
    
    // Disable buttons during request
    const buttons = element.querySelectorAll('button, input[type="submit"]');
    buttons.forEach(btn => {
        btn.disabled = true;
        btn.classList.add('loading');
    });
}

function hideLoadingState(element) {
    const loadingEl = element.querySelector('.htmx-indicator');
    if (loadingEl) {
        loadingEl.style.opacity = '0';
    }
    
    // Re-enable buttons
    const buttons = element.querySelectorAll('button, input[type="submit"]');
    buttons.forEach(btn => {
        btn.disabled = false;
        btn.classList.remove('loading');
    });
}

// Toast notifications using DaisyUI
function showToast(message, type = 'info', duration = 5000) {
    const toast = document.createElement('div');
    const alertClass = getAlertClass(type);
    
    toast.className = `alert ${alertClass} alert-floating animate-slide-in`;
    toast.innerHTML = `
        <div class="flex items-center gap-2">
            ${getAlertIcon(type)}
            <span>${message}</span>
        </div>
        <button class="btn btn-ghost btn-sm" onclick="this.parentElement.remove()">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
        </button>
    `;
    
    document.body.appendChild(toast);
    
    // Auto remove
    setTimeout(() => {
        toast.style.transition = 'opacity 0.3s ease-out, transform 0.3s ease-out';
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 300);
    }, duration);
}

function getAlertClass(type) {
    switch (type) {
        case 'success': return 'alert-success';
        case 'error': return 'alert-error';
        case 'warning': return 'alert-warning';
        case 'info': 
        default: return 'alert-info';
    }
}

function getAlertIcon(type) {
    switch (type) {
        case 'success':
            return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>`;
        case 'error':
            return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>`;
        case 'warning':
            return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.996-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                    </svg>`;
        case 'info':
        default:
            return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>`;
    }
}

// Form validation
function validateForm(form) {
    let isValid = true;
    const requiredFields = form.querySelectorAll('[required]');
    
    requiredFields.forEach(field => {
        const value = field.value.trim();
        const fieldContainer = field.closest('.form-control');
        
        // Remove previous error states
        field.classList.remove('input-error');
        const errorMsg = fieldContainer?.querySelector('.text-error');
        if (errorMsg) errorMsg.remove();
        
        if (!value) {
            isValid = false;
            field.classList.add('input-error');
            
            if (fieldContainer) {
                const errorDiv = document.createElement('div');
                errorDiv.className = 'text-error text-sm mt-1';
                errorDiv.textContent = 'This field is required';
                fieldContainer.appendChild(errorDiv);
            }
        }
        
        // Email validation
        if (field.type === 'email' && value) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                isValid = false;
                field.classList.add('input-error');
                
                if (fieldContainer) {
                    const errorDiv = document.createElement('div');
                    errorDiv.className = 'text-error text-sm mt-1';
                    errorDiv.textContent = 'Please enter a valid email address';
                    fieldContainer.appendChild(errorDiv);
                }
            }
        }
    });
    
    return isValid;
}

// Modal management
function initModals() {
    // Close modals on backdrop click
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('modal') && e.target.classList.contains('modal-open')) {
            closeModal(e.target.id);
        }
    });
    
    // Close modals on escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            const openModal = document.querySelector('.modal.modal-open');
            if (openModal) {
                closeModal(openModal.id);
            }
        }
    });
}

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('modal-open');
        document.body.style.overflow = 'hidden';
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('modal-open');
        document.body.style.overflow = '';
    }
}

// Dropdown management
function initDropdowns() {
    document.addEventListener('click', function(e) {
        const dropdowns = document.querySelectorAll('.dropdown.dropdown-open');
        dropdowns.forEach(dropdown => {
            if (!dropdown.contains(e.target)) {
                dropdown.classList.remove('dropdown-open');
            }
        });
    });
}

// Utility functions for external use
window.whalio = {
    showToast,
    openModal,
    closeModal,
    validateForm,
    showLoadingState,
    hideLoadingState
};

// Progress bar utilities
function showProgressBar(container, progress = 0) {
    const progressBar = document.createElement('div');
    progressBar.className = 'progress progress-primary w-full';
    progressBar.innerHTML = `<div class="progress-bar" style="width: ${progress}%"></div>`;
    container.appendChild(progressBar);
    return progressBar;
}

function updateProgress(progressBar, progress) {
    const bar = progressBar.querySelector('.progress-bar');
    if (bar) {
        bar.style.width = `${progress}%`;
    }
}

// Debounce utility
function debounce(func, wait, immediate) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            timeout = null;
            if (!immediate) func(...args);
        };
        const callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func(...args);
    };
}

// Throttle utility
function throttle(func, limit) {
    let inThrottle;
    return function() {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    }
}

// Export utilities
window.whalio = {
    ...window.whalio,
    showProgressBar,
    updateProgress,
    debounce,
    throttle
};