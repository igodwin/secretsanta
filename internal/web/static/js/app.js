// Application State
const state = {
    participants: [],
    drawResults: null,
    availableNotifiers: []
};

// API Base URL
const API_BASE = window.location.origin;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initializeTabs();
    initializeFormatTabs();
    initializeCreateTab();
    initializeUploadTab();
    initializeDrawTab();
    updateParticipantCount();
    fetchNotificationStatus();
});

// Tab Management
function initializeTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabName = button.dataset.tab;

            // Update buttons
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            // Update content
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabName}-tab`).classList.add('active');

            // Update draw tab when switching to it
            if (tabName === 'draw') {
                updateDrawTab();
            }
        });
    });
}

// Format Tab Management
function initializeFormatTabs() {
    const formatButtons = document.querySelectorAll('.format-tab-btn');
    const formatContents = document.querySelectorAll('.format-content');

    formatButtons.forEach(button => {
        button.addEventListener('click', () => {
            const format = button.dataset.format;

            // Update buttons
            formatButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            // Update content
            formatContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`format-${format}`).classList.add('active');
        });
    });

    // Initialize template download buttons
    const templateButtons = document.querySelectorAll('.download-template');
    templateButtons.forEach(button => {
        button.addEventListener('click', () => {
            const format = button.dataset.format;
            downloadTemplate(format);
        });
    });
}

// Create Tab
function initializeCreateTab() {
    const form = document.getElementById('participant-form');
    const validateBtn = document.getElementById('validate-btn');
    const clearBtn = document.getElementById('clear-btn');

    form.addEventListener('submit', (e) => {
        e.preventDefault();
        addParticipant();
    });

    validateBtn.addEventListener('click', validateParticipants);
    clearBtn.addEventListener('click', clearAllParticipants);
}

function addParticipant() {
    const form = document.getElementById('participant-form');
    const formData = new FormData(form);

    const contactInfo = formData.get('contact_info')
        .split(',')
        .map(s => s.trim())
        .filter(s => s);

    const exclusions = formData.get('exclusions')
        .split(',')
        .map(s => s.trim())
        .filter(s => s);

    const participant = {
        name: formData.get('name').trim(),
        notification_type: formData.get('notification_type'),
        contact_info: contactInfo,
        exclusions: exclusions
    };

    // Check for duplicate names
    if (state.participants.some(p => p.name === participant.name)) {
        showToast('Participant with this name already exists!', 'error');
        return;
    }

    state.participants.push(participant);
    renderParticipants();
    updateParticipantCount();
    form.reset();
    showToast('Participant added successfully!', 'success');
}

function renderParticipants() {
    const container = document.getElementById('participants-container');

    if (state.participants.length === 0) {
        container.innerHTML = '<p style="color: #6c757d; text-align: center; padding: 20px;">No participants added yet</p>';
        return;
    }

    container.innerHTML = state.participants.map((p, index) => `
        <div class="participant-card">
            <div class="participant-info">
                <strong>${escapeHtml(p.name)}</strong>
                <small>
                    ${escapeHtml(p.notification_type)} • ${escapeHtml(p.contact_info.join(', '))}
                    ${p.exclusions.length > 0 ? ` • Excludes: ${escapeHtml(p.exclusions.join(', '))}` : ''}
                </small>
            </div>
            <button onclick="removeParticipant(${index})">Remove</button>
        </div>
    `).join('');
}

function removeParticipant(index) {
    state.participants.splice(index, 1);
    renderParticipants();
    updateParticipantCount();
    showToast('Participant removed', 'success');
}

function clearAllParticipants() {
    if (state.participants.length === 0) return;

    if (confirm('Are you sure you want to clear all participants?')) {
        state.participants = [];
        renderParticipants();
        updateParticipantCount();
        showToast('All participants cleared', 'success');
    }
}

function updateParticipantCount() {
    document.getElementById('participant-count').textContent = state.participants.length;
    document.getElementById('draw-participant-count').textContent = state.participants.length;

    // Enable/disable draw button
    const drawBtn = document.getElementById('run-draw-btn');
    drawBtn.disabled = state.participants.length < 2;
}

async function validateParticipants() {
    if (state.participants.length < 2) {
        showToast('Need at least 2 participants to validate', 'warning');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/api/validate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(state.participants)
        });

        const result = await response.json();
        showValidationModal(result);
    } catch (error) {
        showToast('Validation failed: ' + error.message, 'error');
    }
}

// Upload Tab
function initializeUploadTab() {
    const uploadArea = document.getElementById('upload-area');
    const fileInput = document.getElementById('file-input');
    const browseBtn = document.getElementById('browse-btn');

    browseBtn.addEventListener('click', () => fileInput.click());
    fileInput.addEventListener('change', handleFileSelect);

    // Drag and drop
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.classList.add('drag-over');
    });

    uploadArea.addEventListener('dragleave', () => {
        uploadArea.classList.remove('drag-over');
    });

    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.classList.remove('drag-over');
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFile(files[0]);
        }
    });
}

function handleFileSelect(e) {
    const file = e.target.files[0];
    if (file) {
        handleFile(file);
    }
}

async function handleFile(file) {
    // Check if file extension is supported
    const validExtensions = ['.json', '.yaml', '.yml', '.toml', '.csv', '.tsv'];
    const fileExt = file.name.substring(file.name.lastIndexOf('.')).toLowerCase();

    if (!validExtensions.includes(fileExt)) {
        showToast('Please upload a supported file format (JSON, YAML, TOML, CSV, TSV)', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch(`${API_BASE}/api/upload`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        const result = await response.json();
        state.participants = result.participants;
        renderParticipants();
        updateParticipantCount();

        showToast(`Loaded ${result.participants.length} participants from ${result.format.toUpperCase()}`, 'success');

        // Show validation results
        if (result.validation) {
            setTimeout(() => showValidationModal(result.validation), 500);
        }

        // Switch to create tab to show loaded participants
        document.querySelector('[data-tab="create"]').click();
    } catch (error) {
        showToast('Upload failed: ' + error.message, 'error');
    }
}

async function downloadTemplate(format) {
    try {
        const response = await fetch(`${API_BASE}/api/template?format=${format}`);

        if (!response.ok) {
            throw new Error('Failed to download template');
        }

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;

        // Get filename from Content-Disposition header or use default
        const contentDisposition = response.headers.get('Content-Disposition');
        let filename = `secretsanta-template.${format}`;
        if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename=(.+)/);
            if (filenameMatch) {
                filename = filenameMatch[1];
            }
        }

        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);
        showToast(`Template downloaded (${format.toUpperCase()})`, 'success');
    } catch (error) {
        showToast('Failed to download template: ' + error.message, 'error');
    }
}

// Draw Tab
function initializeDrawTab() {
    const runDrawBtn = document.getElementById('run-draw-btn');
    const exportBtn = document.getElementById('export-btn');
    const newDrawBtn = document.getElementById('new-draw-btn');

    runDrawBtn.addEventListener('click', runDraw);
    exportBtn.addEventListener('click', exportResults);
    newDrawBtn.addEventListener('click', resetDraw);
}

function updateDrawTab() {
    const validationStatus = document.getElementById('validation-status');

    if (state.participants.length < 2) {
        validationStatus.innerHTML = '<span style="color: var(--danger-color);">⚠️ Need at least 2 participants</span>';
        return;
    }

    validationStatus.innerHTML = '<span style="color: var(--success-color);">✓ Ready to draw</span>';
}

async function runDraw() {
    if (state.participants.length < 2) {
        showToast('Need at least 2 participants', 'error');
        return;
    }

    const runDrawBtn = document.getElementById('run-draw-btn');
    runDrawBtn.disabled = true;
    runDrawBtn.textContent = 'Drawing...';

    try {
        // Get archive email if provided
        const archiveEmail = document.getElementById('archive-email').value.trim();

        // Build request payload
        const requestBody = {
            participants: state.participants
        };

        // Add archive email if provided
        if (archiveEmail) {
            requestBody.archive_email = archiveEmail;
        }

        const response = await fetch(`${API_BASE}/api/draw`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(requestBody)
        });

        const result = await response.json();

        if (!result.success) {
            throw new Error(result.error);
        }

        state.drawResults = result.participants;
        displayResults();

        if (archiveEmail) {
            showToast(`Draw completed! Archive sent to ${archiveEmail}`, 'success');
        } else {
            showToast('Draw completed successfully!', 'success');
        }
    } catch (error) {
        showToast('Draw failed: ' + error.message, 'error');
        runDrawBtn.disabled = false;
        runDrawBtn.textContent = 'Run Draw';
    }
}

function displayResults() {
    const resultsSection = document.getElementById('results-section');
    const resultsContainer = document.getElementById('results-container');
    const runDrawBtn = document.getElementById('run-draw-btn');

    resultsSection.style.display = 'block';
    runDrawBtn.style.display = 'none';

    // Show the warning instead of the actual results
    resultsContainer.innerHTML = `
        <div class="papa-elf-warning">
            <div class="elf-message">
                <h3>🎅 HOLD IT RIGHT THERE, BUDDY! 🎄</h3>
                <p style="font-size: 1.1em; margin: 20px 0;">
                    Papa Elf here! Now listen up, sport - those Secret Santa assignments are
                    <strong>TOP SECRET</strong>, just like the Claus family cookie recipe!
                </p>
                <p style="font-size: 1em; color: #666; margin: 15px 0;">
                    Peeking at these results might spoil the MAGIC and WONDER of the holiday season!
                    You could ruin surprises faster than putting maple syrup on spaghetti
                    (which is delicious, by the way).
                </p>
                <p style="font-size: 1em; color: #d32f2f; margin: 15px 0; font-weight: bold;">
                    ⚠️ Are you ABSOLUTELY, POSITIVELY, 100% SURE you want to see who's giving to whom? ⚠️
                </p>
                <p style="font-size: 0.9em; color: #888; margin: 10px 0;">
                    (Remember: With great power comes great responsibility to keep secrets!)
                </p>
            </div>
            <button id="reveal-results-btn" class="btn btn-warning btn-large" style="margin: 20px 0;">
                Yes, I Understand the Consequences - Reveal Results
            </button>
            <p style="font-size: 0.85em; color: #999; margin-top: 10px;">
                Don't say Papa Elf didn't warn you! 🤶
            </p>
        </div>
    `;

    // Add event listener for reveal button
    document.getElementById('reveal-results-btn').addEventListener('click', confirmAndRevealResults);
}

function confirmAndRevealResults() {
    const confirmation = confirm(
        "🎄 FINAL WARNING FROM PAPA ELF! 🎄\n\n" +
        "You're about to reveal ALL the Secret Santa assignments!\n\n" +
        "Once you see them, you can't unsee them - it's like trying to forget " +
        "you saw Santa eating cookies at 3am!\n\n" +
        "The magic of surprise will be GONE! KAPUT! FINITO!\n\n" +
        "Are you really, truly, cross-your-heart sure you want to proceed?\n\n" +
        "Click OK if you're ready to shoulder this tremendous responsibility,\n" +
        "or Cancel if you'd rather keep the Christmas spirit alive!"
    );

    if (confirmation) {
        revealActualResults();
    }
}

function revealActualResults() {
    const resultsContainer = document.getElementById('results-container');

    resultsContainer.innerHTML = `
        <div class="revealed-warning" style="background: #fff3cd; padding: 15px; border-radius: 8px; margin-bottom: 20px; border: 2px solid #ffc107;">
            <p style="margin: 0; color: #856404; font-size: 0.95em;">
                🎁 <strong>Papa Elf says:</strong> You've been entrusted with the sacred knowledge!
                Guard it well, like the elves guard the Naughty & Nice List!
            </p>
        </div>
        ${state.drawResults.map(p => `
            <div class="result-card">
                <div class="giver">${escapeHtml(p.name)}</div>
                <div class="arrow">→</div>
                <div class="recipient">${escapeHtml(p.recipient)}</div>
            </div>
        `).join('')}
    `;
}

function exportResults() {
    if (!state.drawResults) return;

    const blob = new Blob([JSON.stringify(state.drawResults, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'secretsanta-results.json';
    a.click();
    URL.revokeObjectURL(url);
    showToast('Results exported', 'success');
}

function resetDraw() {
    state.drawResults = null;
    document.getElementById('results-section').style.display = 'none';
    document.getElementById('run-draw-btn').style.display = 'block';
    document.getElementById('run-draw-btn').disabled = false;
    document.getElementById('run-draw-btn').textContent = 'Run Draw';
}

// Validation Modal
function showValidationModal(result) {
    const modal = document.getElementById('validation-modal');
    const resultsDiv = document.getElementById('validation-results');

    let html = '';

    if (result.valid) {
        html += `<div class="validation-success">
            <h3>✓ Configuration is Valid</h3>
            <p><strong>Total Participants:</strong> ${result.total_participants}</p>
            <p><strong>Minimum Compatibility:</strong> ${result.min_compatibility}</p>
            <p><strong>Average Compatibility:</strong> ${result.avg_compatibility.toFixed(1)}</p>
        </div>`;
    } else {
        html += `<div class="validation-error">
            <h3>✗ Configuration is Invalid</h3>
            <h4>Errors:</h4>
            <ul>${result.errors.map(e => `<li>${escapeHtml(e)}</li>`).join('')}</ul>
        </div>`;
    }

    if (result.warnings && result.warnings.length > 0) {
        html += `<div class="validation-warning">
            <h4>⚠️ Warnings:</h4>
            <ul>${result.warnings.map(w => `<li>${escapeHtml(w)}</li>`).join('')}</ul>
        </div>`;
    }

    resultsDiv.innerHTML = html;
    modal.style.display = 'block';
}

// Modal close handlers
document.querySelector('.close').addEventListener('click', () => {
    document.getElementById('validation-modal').style.display = 'none';
});

window.addEventListener('click', (e) => {
    const modal = document.getElementById('validation-modal');
    if (e.target === modal) {
        modal.style.display = 'none';
    }
});

// Toast Notifications
function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type}`;
    toast.style.display = 'block';

    setTimeout(() => {
        toast.style.display = 'none';
    }, 3000);
}

// Utility Functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Notification Status
async function fetchNotificationStatus() {
    const statusContainer = document.getElementById('notification-types');

    try {
        const response = await fetch(`${API_BASE}/api/status`);
        if (!response.ok) {
            throw new Error('Failed to fetch status');
        }

        const status = await response.json();
        displayNotificationStatus(status);
    } catch (error) {
        console.error('Error fetching notification status:', error);
        statusContainer.innerHTML = '<span class="status-error">Unknown</span>';
    }
}

function displayNotificationStatus(status) {
    const container = document.getElementById('notification-types');

    if (!status.available || status.available.length === 0) {
        container.innerHTML = '<span class="status-type">None configured</span>';
        return;
    }

    // Store available notifiers in state
    state.availableNotifiers = status.available;

    // Update the notification type dropdown
    updateNotificationTypeDropdown();

    // Create badges for each notification type
    const badges = status.available.map(notifier => {
        const type = notifier.type || notifier; // Support both object and string format
        const icon = getNotificationIcon(type);
        const accountInfo = notifier.accounts && notifier.accounts.length > 0
            ? ` (${notifier.accounts.length} account${notifier.accounts.length > 1 ? 's' : ''})`
            : '';
        return `<span class="status-type" title="${type}${accountInfo}">${icon} ${type}</span>`;
    }).join('');

    let html = badges;

    // Add notifier service status if configured
    if (status.using_notifier) {
        const healthIcon = status.notifier_healthy ? '✓' : '✗';
        const healthClass = status.notifier_healthy ? 'healthy' : 'unhealthy';
        html += ` <span class="notifier-status ${healthClass}" title="External notifier service: ${status.notifier_status || 'unknown'}">(${healthIcon} notifier)</span>`;
    } else if (status.smtp_configured) {
        html += ` <span class="notifier-status healthy" title="Using built-in SMTP">(✓ SMTP)</span>`;
    }

    container.innerHTML = html;
}

function updateNotificationTypeDropdown() {
    const dropdown = document.getElementById('notification-type');
    if (!dropdown) return;

    // Clear existing options
    dropdown.innerHTML = '';

    // Populate with available notifier types
    state.availableNotifiers.forEach(notifier => {
        const type = notifier.type || notifier; // Support both object and string format

        // If this notifier has multiple accounts, create an option for each
        if (notifier.accounts && notifier.accounts.length > 0) {
            notifier.accounts.forEach(account => {
                const option = document.createElement('option');
                // Format: type:account (e.g., "email:notify")
                option.value = `${type}:${account}`;

                // Create readable label
                const typeLabel = type.charAt(0).toUpperCase() + type.slice(1);
                const defaultMarker = account === notifier.default_account ? ' (default)' : '';
                option.textContent = `${typeLabel} - ${account}${defaultMarker}`;

                dropdown.appendChild(option);
            });
        } else {
            // No accounts, just add the type
            const option = document.createElement('option');
            option.value = type;
            const typeLabel = type.charAt(0).toUpperCase() + type.slice(1);
            option.textContent = typeLabel;
            dropdown.appendChild(option);
        }
    });

    // If no notifiers available, add a default option
    if (state.availableNotifiers.length === 0) {
        const option = document.createElement('option');
        option.value = 'stdout';
        option.textContent = 'Console (stdout)';
        dropdown.appendChild(option);
    }
}

function getNotificationIcon(type) {
    const icons = {
        'email': '📧',
        'slack': '💬',
        'ntfy': '🔔',
        'stdout': '💻'
    };
    return icons[type.toLowerCase()] || '📬';
}
