document.addEventListener('DOMContentLoaded', () => {
    // Your existing code goes inside here safely
    const formElement = document.getElementById('tm-form');
    
    if (formElement) {
        formElement.addEventListener('submit', async (e) => {
            e.preventDefault();

            const initialTape = document.getElementById('initial-tape').value;
            const machineJson = document.getElementById('machine-json').value;
            const outputContainer = document.getElementById('output-container');

            outputContainer.innerHTML = '<p class="text-blue-400">Running simulation...</p>';

            try {
                const response = await fetch('/api/run', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        initial_tape: initialTape,
                        machine_json: machineJson
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

        let statusColor = 'text-green-300';
        if (step.Status === 'Rejected' || step.Status === 'Timeout') statusColor = 'text-red-400';
        else if (step.Status === 'Accepted') statusColor = 'text-yellow-400 font-bold';

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