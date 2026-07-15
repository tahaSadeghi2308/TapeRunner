async function loadMachines() {
    const machineSelect = document.getElementById('machine-select');
    if (!machineSelect) return;

    try {
        const response = await fetch('/api/machines');
        if (!response.ok) {
            throw new Error('Failed to load machines');
        }

        const machines = await response.json();
        machineSelect.innerHTML = '';

        if (machines.length === 0) {
            const option = document.createElement('option');
            option.value = '';
            option.textContent = 'No machines available';
            machineSelect.appendChild(option);
            return;
        }

        machines.forEach(machine => {
            const option = document.createElement('option');
            option.value = machine;
            option.textContent = machine;
            machineSelect.appendChild(option);
        });
    } catch (error) {
        const option = document.createElement('option');
        option.value = '';
        option.textContent = 'Error loading machines';
        machineSelect.appendChild(option);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    loadMachines();

    const formElement = document.getElementById('tm-form');
    
    if (formElement) {
        formElement.addEventListener('submit', async (e) => {
            e.preventDefault();

            const initialTape = document.getElementById('initial-tape').value;
            const machineSelect = document.getElementById('machine-select');
            const machineName = machineSelect.value;
            const outputContainer = document.getElementById('output-container');

            if (!machineName) {
                outputContainer.innerHTML = '<p class="text-red-500 font-bold">Error: Please select a machine</p>';
                return;
            }

            outputContainer.innerHTML = '<p class="text-blue-400">Running simulation...</p>';

            try {
                const response = await fetch('/api/run', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        machine_name: machineName,
                        initial_tape: initialTape
                    })
                });

                if (!response.ok) {
                    const errText = await response.text();
                    throw new Error(errText);
                }

                const steps = await response.json();
                renderSteps(steps, outputContainer);

            } catch (error) {
                outputContainer.innerHTML = `<p class="text-red-500 font-bold">Error: ${error.message}</p>`;
            }
        });
    }
});

function renderSteps(steps, container) {
    container.innerHTML = ''; 

    steps.forEach((step, index) => {
        const stepDiv = document.createElement('div');
        stepDiv.className = 'mb-6 p-3 bg-gray-800 rounded shadow';

        let statusColor = 'text-yellow-400';
        if (step.Status === 'Rejected' || step.Status === 'Timeout') statusColor = 'text-red-400';
        else if (step.Status === 'Accepted') statusColor = 'text-green-300 font-bold';

        const headerHtml = `
            <div class="flex justify-between mb-3 border-b border-gray-700 pb-2">
                <span class="text-blue-300 font-bold">Step: ${index} | State: <span class="text-white">${step.CurrentState}</span></span>
                <span class="${statusColor}">Status: ${step.Status}</span>
            </div>
        `;

        let tapeHtml = '<div class="flex space-x-1 overflow-x-auto pb-4">';
        for (let i = 0; i < step.TapeContent.length; i++) {
            const char = step.TapeContent[i];
            const isHead = i === step.HeadPosition;
            
            const cellBg = isHead ? 'bg-blue-900 border-blue-500' : 'bg-gray-700 border-gray-500';
            const headIndicator = isHead ? '<div class="absolute -bottom-4 text-xs text-blue-400 font-bold">▲</div>' : '';

            tapeHtml += `
                <div class="tape-cell flex-shrink-0 flex items-center justify-center border text-lg relative ${cellBg}">
                    ${char}
                    ${headIndicator}
                </div>
            `;
        }
        tapeHtml += '</div>';

        stepDiv.innerHTML = headerHtml + tapeHtml;
        container.appendChild(stepDiv);
    });
}
