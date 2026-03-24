/**
 * Chain Lens - Bitcoin Transaction Visualizer
 * A user-friendly interface for understanding Bitcoin transactions
 */

// ===== Demo Data =====
const DEMO_TRANSACTIONS = {
    segwit: {
        ok: true,
        network: "mainnet",
        segwit: true,
        txid: "e4c9d8f5a2b3c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6",
        wtxid: "f5d9c8e4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6",
        version: 2,
        locktime: 0,
        size_bytes: 225,
        weight: 573,
        vbytes: 144,
        total_input_sats: 100000,
        total_output_sats: 98500,
        fee_sats: 1500,
        fee_rate_sat_vb: 10.42,
        rbf_signaling: false,
        locktime_type: "none",
        locktime_value: 0,
        segwit_savings: {
            witness_bytes: 108,
            non_witness_bytes: 117,
            total_bytes: 225,
            weight_actual: 573,
            weight_if_legacy: 900,
            savings_pct: 36.33
        },
        vin: [
            {
                txid: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
                vout: 0,
                sequence: 4294967295,
                script_sig_hex: "",
                script_asm: "",
                witness: ["3044022...", "02a1b2c3..."],
                script_type: "p2wpkh",
                address: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
                prevout: {
                    value_sats: 100000,
                    script_pubkey_hex: "0014751e76e8199196d454941c45d1b3a323f1433bd6"
                },
                relative_timelock: { enabled: false }
            }
        ],
        vout: [
            {
                n: 0,
                value_sats: 50000,
                script_pubkey_hex: "0014d85c2b71d0060b09c9886aeb815e50991dda124d",
                script_asm: "OP_0 OP_PUSHBYTES_20 d85c2b71d0060b09c9886aeb815e50991dda124d",
                script_type: "p2wpkh",
                address: "bc1qmzwtfh8gqctsnye5xhtwr90gxjwm6pyng7cdh4"
            },
            {
                n: 1,
                value_sats: 48500,
                script_pubkey_hex: "0014751e76e8199196d454941c45d1b3a323f1433bd6",
                script_asm: "OP_0 OP_PUSHBYTES_20 751e76e8199196d454941c45d1b3a323f1433bd6",
                script_type: "p2wpkh",
                address: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
            }
        ],
        warnings: []
    },
    legacy: {
        ok: true,
        network: "mainnet",
        segwit: false,
        txid: "b4c9d8e5a2b3c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6",
        wtxid: null,
        version: 1,
        locktime: 500000,
        size_bytes: 226,
        weight: 904,
        vbytes: 226,
        total_input_sats: 500000,
        total_output_sats: 485000,
        fee_sats: 15000,
        fee_rate_sat_vb: 66.37,
        rbf_signaling: true,
        locktime_type: "block_height",
        locktime_value: 500000,
        segwit_savings: null,
        vin: [
            {
                txid: "c1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
                vout: 1,
                sequence: 4294967293,
                script_sig_hex: "483045022100...",
                script_asm: "OP_PUSHBYTES_72 3045022100... OP_PUSHBYTES_33 02a1b2c3...",
                witness: [],
                script_type: "p2pkh",
                address: "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
                prevout: {
                    value_sats: 500000,
                    script_pubkey_hex: "76a91477bff20c60e522dfaa"
                },
                relative_timelock: { enabled: false }
            }
        ],
        vout: [
            {
                n: 0,
                value_sats: 400000,
                script_pubkey_hex: "76a914...",
                script_asm: "OP_DUP OP_HASH160 OP_PUSHBYTES_20 ... OP_EQUALVERIFY OP_CHECKSIG",
                script_type: "p2pkh",
                address: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
            },
            {
                n: 1,
                value_sats: 85000,
                script_pubkey_hex: "76a91477bff...",
                script_asm: "OP_DUP OP_HASH160 OP_PUSHBYTES_20 ... OP_EQUALVERIFY OP_CHECKSIG",
                script_type: "p2pkh",
                address: "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"
            }
        ],
        warnings: [
            { code: "HIGH_FEE" },
            { code: "RBF_SIGNALING" }
        ]
    },
    opreturn: {
        ok: true,
        network: "mainnet",
        segwit: true,
        txid: "d4c9d8e5a2b3c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6",
        wtxid: "e5d9c8e4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6",
        version: 2,
        locktime: 0,
        size_bytes: 250,
        weight: 610,
        vbytes: 153,
        total_input_sats: 50000,
        total_output_sats: 48000,
        fee_sats: 2000,
        fee_rate_sat_vb: 13.07,
        rbf_signaling: false,
        locktime_type: "none",
        locktime_value: 0,
        segwit_savings: {
            witness_bytes: 108,
            non_witness_bytes: 142,
            total_bytes: 250,
            weight_actual: 610,
            weight_if_legacy: 1000,
            savings_pct: 39.0
        },
        vin: [
            {
                txid: "f1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
                vout: 0,
                sequence: 4294967295,
                script_sig_hex: "",
                script_asm: "",
                witness: ["3044022...", "02a1b2c3..."],
                script_type: "p2wpkh",
                address: "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
                prevout: {
                    value_sats: 50000,
                    script_pubkey_hex: "0014e8df018c7e326cc253fabd8ab2be1e8f3c5dbafc"
                },
                relative_timelock: { enabled: false }
            }
        ],
        vout: [
            {
                n: 0,
                value_sats: 48000,
                script_pubkey_hex: "0014e8df018c7e326cc253fabd8ab2be1e8f3c5dbafc",
                script_asm: "OP_0 OP_PUSHBYTES_20 e8df018c7e326cc253fabd8ab2be1e8f3c5dbafc",
                script_type: "p2wpkh",
                address: "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
            },
            {
                n: 1,
                value_sats: 0,
                script_pubkey_hex: "6a1468656c6c6f2c20626974636f696e21",
                script_asm: "OP_RETURN OP_PUSHBYTES_20 68656c6c6f2c20626974636f696e21",
                script_type: "op_return",
                address: null,
                op_return_data_hex: "68656c6c6f2c20626974636f696e21",
                op_return_data_utf8: "hello, bitcoin!",
                op_return_protocol: "text"
            }
        ],
        warnings: [
            { code: "DUST_OUTPUT" }
        ]
    }
};

// ===== DOM Elements =====
const elements = {
    // Panels
    inputPanel: document.getElementById('input-panel'),
    resultsPanel: document.getElementById('results-panel'),
    errorPanel: document.getElementById('error-panel'),

    // Input
    jsonInput: document.getElementById('json-input'),
    btnAnalyze: document.getElementById('btn-analyze'),
    btnBack: document.getElementById('btn-back'),
    btnErrorBack: document.getElementById('btn-error-back'),

    // Story
    storyWhoMain: document.getElementById('story-who-main'),
    storyWhoDetail: document.getElementById('story-who-detail'),
    storyCostMain: document.getElementById('story-cost-main'),
    storyCostDetail: document.getElementById('story-cost-detail'),
    feeMeter: document.getElementById('fee-meter'),
    feeMeterFill: document.getElementById('fee-meter-fill'),
    feeMeterLabel: document.getElementById('fee-meter-label'),
    storyRiskCard: document.getElementById('story-risk-card'),
    storyRiskMain: document.getElementById('story-risk-main'),
    riskList: document.getElementById('risk-list'),

    // Flow
    inputsList: document.getElementById('inputs-list'),
    outputsList: document.getElementById('outputs-list'),
    totalIn: document.getElementById('total-in'),
    totalOut: document.getElementById('total-out'),
    feeValue: document.getElementById('fee-value'),
    balanceOutputs: document.getElementById('balance-outputs'),
    balanceOutputsLabel: document.getElementById('balance-outputs-label'),
    balanceFee: document.getElementById('balance-fee'),
    balanceFeeLabel: document.getElementById('balance-fee-label'),

    // Tech
    techTxid: document.getElementById('tech-txid'),
    techWtxid: document.getElementById('tech-wtxid'),
    wtxidCard: document.getElementById('wtxid-card'),
    techType: document.getElementById('tech-type'),
    techTypeHint: document.getElementById('tech-type-hint'),
    techSize: document.getElementById('tech-size'),
    techVbytes: document.getElementById('tech-vbytes'),
    techWeight: document.getElementById('tech-weight'),
    techFeerate: document.getElementById('tech-feerate'),
    feeRateHint: document.getElementById('fee-rate-hint'),
    techVersion: document.getElementById('tech-version'),
    techLocktime: document.getElementById('tech-locktime'),
    locktimeHint: document.getElementById('locktime-hint'),
    techRbf: document.getElementById('tech-rbf'),

    // SegWit Savings
    segwitSavings: document.getElementById('segwit-savings'),
    savingsActual: document.getElementById('savings-actual'),
    savingsActualLabel: document.getElementById('savings-actual-label'),
    savingsLegacy: document.getElementById('savings-legacy'),
    savingsLegacyLabel: document.getElementById('savings-legacy-label'),
    savingsText: document.getElementById('savings-text'),

    // Other
    rawJson: document.getElementById('raw-json'),
    tooltip: document.getElementById('tooltip'),
    errorCode: document.getElementById('error-code'),
    errorMessage: document.getElementById('error-message')
};

// ===== Utility Functions =====

/**
 * Format satoshis into a human-readable string
 */
function formatSats(sats) {
    if (sats >= 100000000) {
        return (sats / 100000000).toFixed(8).replace(/\.?0+$/, '') + ' BTC';
    } else if (sats >= 1000000) {
        return (sats / 1000000).toFixed(2) + ' mBTC';
    } else if (sats >= 1000) {
        return sats.toLocaleString() + ' sats';
    }
    return sats + ' sats';
}

/**
 * Format satoshis with always showing sats
 */
function formatSatsWithUnit(sats) {
    if (sats >= 100000000) {
        const btc = (sats / 100000000).toFixed(8).replace(/\.?0+$/, '');
        return `${btc} BTC (${sats.toLocaleString()} sats)`;
    }
    return sats.toLocaleString() + ' sats';
}

/**
 * Truncate an address for display
 */
function truncateAddress(address) {
    if (!address || address.length <= 20) return address || 'Unknown';
    return address.slice(0, 10) + '...' + address.slice(-8);
}

/**
 * Get human-readable script type name
 */
function getScriptTypeName(type) {
    const names = {
        'p2pkh': 'P2PKH (Legacy)',
        'p2sh': 'P2SH (Script Hash)',
        'p2sh-p2wpkh': 'P2SH-P2WPKH (Wrapped SegWit)',
        'p2sh-p2wsh': 'P2SH-P2WSH (Wrapped SegWit)',
        'p2wpkh': 'P2WPKH (Native SegWit)',
        'p2wsh': 'P2WSH (SegWit Script)',
        'p2tr': 'P2TR (Taproot)',
        'p2tr_keypath': 'P2TR Key Path',
        'p2tr_scriptpath': 'P2TR Script Path',
        'op_return': 'OP_RETURN (Data)',
        'unknown': 'Unknown'
    };
    return names[type] || type?.toUpperCase() || 'Unknown';
}

/**
 * Format warning code into human-readable text
 */
function formatWarning(code) {
    const warnings = {
        'HIGH_FEE': 'High fee — this transaction paid more than typical',
        'DUST_OUTPUT': 'Dust output — very small value that may be uneconomical to spend',
        'UNKNOWN_OUTPUT_SCRIPT': 'Unknown script — non-standard output detected',
        'RBF_SIGNALING': 'RBF enabled — sender can replace this with a higher fee before confirmation'
    };
    return warnings[code] || code;
}

/**
 * Get fee rate classification
 */
function getFeeRateHint(feeRate) {
    if (feeRate < 2) return 'Very low — may take days to confirm';
    if (feeRate < 5) return 'Low priority';
    if (feeRate < 15) return 'Normal priority';
    if (feeRate < 50) return 'High priority';
    return 'Very high — premium for fast confirmation';
}

// ===== Tab Switching =====
document.querySelectorAll('.input-tab').forEach(tab => {
    tab.addEventListener('click', () => {
        document.querySelectorAll('.input-tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));

        tab.classList.add('active');
        const panelId = tab.dataset.tab + '-panel';
        document.getElementById(panelId).classList.add('active');
    });
});

// ===== Demo Buttons =====
document.querySelectorAll('.demo-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const demoType = btn.dataset.demo;
        const demoData = DEMO_TRANSACTIONS[demoType];
        if (demoData) {
            displayResults(demoData);
        }
    });
});

// ===== Analyze Button =====
elements.btnAnalyze.addEventListener('click', async () => {
    const jsonText = elements.jsonInput.value.trim();

    if (!jsonText) {
        showError('EMPTY_INPUT', 'Please paste transaction JSON data');
        return;
    }

    try {
        const data = JSON.parse(jsonText);

        // If this is a fixture (has raw_tx), call the API to analyze it
        if (data.raw_tx && data.prevouts) {
            try {
                const response = await fetch('/api/analyze', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: jsonText
                });
                const result = await response.json();
                if (!result.ok) {
                    showError(result.error?.code || 'API_ERROR', result.error?.message || 'Analysis failed');
                    return;
                }
                displayResults(result);
                return;
            } catch (fetchErr) {
                showError('API_ERROR', 'Failed to call analysis API: ' + fetchErr.message);
                return;
            }
        }

        // Otherwise treat as pre-analyzed JSON output
        if (!data.ok) {
            if (data.error) {
                showError(data.error.code, data.error.message);
            } else {
                showError('PARSE_ERROR', 'Transaction data indicates an error occurred');
            }
            return;
        }

        displayResults(data);
    } catch (e) {
        showError('INVALID_JSON', 'Invalid JSON format: ' + e.message);
    }
});

// ===== Back Buttons =====
elements.btnBack.addEventListener('click', showInputPanel);
elements.btnErrorBack.addEventListener('click', showInputPanel);

// ===== Panel Navigation =====
function showInputPanel() {
    elements.inputPanel.style.display = 'flex';
    elements.resultsPanel.style.display = 'none';
    elements.errorPanel.style.display = 'none';
    const blockPanel = document.getElementById('block-results-panel');
    if (blockPanel) blockPanel.style.display = 'none';
}

function showResultsPanel() {
    elements.inputPanel.style.display = 'none';
    elements.resultsPanel.style.display = 'block';
    elements.errorPanel.style.display = 'none';
    const blockPanel = document.getElementById('block-results-panel');
    if (blockPanel) blockPanel.style.display = 'none';
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

function showBlockResultsPanel() {
    elements.inputPanel.style.display = 'none';
    elements.resultsPanel.style.display = 'none';
    elements.errorPanel.style.display = 'none';
    document.getElementById('block-results-panel').style.display = 'block';
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

function showError(code, message) {
    elements.inputPanel.style.display = 'none';
    elements.resultsPanel.style.display = 'none';
    elements.errorPanel.style.display = 'flex';
    const blockPanel = document.getElementById('block-results-panel');
    if (blockPanel) blockPanel.style.display = 'none';
    elements.errorCode.textContent = code;
    elements.errorMessage.textContent = message;
}

// ===== Display Results =====
function displayResults(data) {
    // Update story cards
    updateStorySection(data);

    // Update flow visualization
    updateFlowSection(data);

    // Update tech details
    updateTechSection(data);

    // Update raw JSON
    elements.rawJson.textContent = JSON.stringify(data, null, 2);

    // Show results panel
    showResultsPanel();
}

// ===== Story Section =====
function updateStorySection(data) {
    const inputCount = data.vin.length;
    const outputCount = data.vout.length;
    const totalIn = data.total_input_sats;
    const totalOut = data.total_output_sats;
    const fee = data.fee_sats;

    // Who paid whom
    const recipientCount = data.vout.filter(v => v.script_type !== 'op_return').length;
    const hasOpReturn = data.vout.some(v => v.script_type === 'op_return');

    elements.storyWhoMain.textContent = `${inputCount} source${inputCount > 1 ? 's' : ''} → ${recipientCount} recipient${recipientCount > 1 ? 's' : ''}`;

    let whoDetail = `Someone spent ${formatSats(totalIn)} from ${inputCount} input${inputCount > 1 ? 's' : ''} `;
    whoDetail += `and sent it to ${recipientCount} address${recipientCount > 1 ? 'es' : ''}.`;
    if (hasOpReturn) {
        whoDetail += ' Also embedded some data on the blockchain.';
    }
    elements.storyWhoDetail.textContent = whoDetail;

    // What did it cost
    elements.storyCostMain.textContent = formatSatsWithUnit(fee) + ' in fees';

    const feePercent = totalIn > 0 ? ((fee / totalIn) * 100).toFixed(2) : 0;
    elements.storyCostDetail.textContent = `That's ${feePercent}% of the total value, paying ${data.fee_rate_sat_vb.toFixed(2)} sat/vB for ${data.vbytes} vbytes.`;

    // Fee meter (scale: 0-100 sat/vB maps to 0-100%)
    const meterPct = Math.min(data.fee_rate_sat_vb, 100);
    elements.feeMeterFill.style.width = meterPct + '%';
    elements.feeMeterLabel.textContent = getFeeRateHint(data.fee_rate_sat_vb);

    // Anything risky
    const warnings = data.warnings || [];
    if (warnings.length === 0) {
        elements.storyRiskCard.classList.remove('has-warnings');
        elements.storyRiskMain.textContent = 'All clear!';
        elements.riskList.innerHTML = '<li style="list-style:none;padding-left:0;">✅ No warnings detected</li>';
        elements.riskList.querySelector('li').style.setProperty('--before', 'none');
    } else {
        elements.storyRiskCard.classList.add('has-warnings');
        elements.storyRiskMain.textContent = `${warnings.length} thing${warnings.length > 1 ? 's' : ''} to watch:`;
        elements.riskList.innerHTML = warnings.map(w =>
            `<li>${formatWarning(w.code)}</li>`
        ).join('');
    }
}

// ===== Flow Section =====
function updateFlowSection(data) {
    // Clear existing items
    elements.inputsList.innerHTML = '';
    elements.outputsList.innerHTML = '';

    // Render inputs
    data.vin.forEach((input, i) => {
        const value = input.prevout?.value_sats || 0;
        const div = document.createElement('div');
        div.className = 'flow-item input-item';

        // Build sequence/RBF info
        let seqHtml = '';
        if (input.sequence !== undefined) {
            const seqHex = '0x' + (input.sequence >>> 0).toString(16).toUpperCase().padStart(8, '0');
            const isRbf = input.sequence < 0xFFFFFFFE;
            seqHtml = `<div class="input-detail"><span class="detail-label">Sequence:</span> <span class="monospace">${seqHex}</span>`;
            if (isRbf) seqHtml += ` <span class="rbf-tag">RBF</span>`;
            seqHtml += `</div>`;
        }

        // Relative timelock info
        let timelockHtml = '';
        if (input.relative_timelock && input.relative_timelock.enabled) {
            const rtl = input.relative_timelock;
            const typeStr = rtl.type === 'blocks' ? `${rtl.value} blocks` : `${rtl.value} × 512 seconds`;
            timelockHtml = `<div class="input-detail timelock-detail">
                <span class="detail-label">⏱ Relative Timelock:</span> ${typeStr}
            </div>`;
        }

        // Witness data preview
        let witnessHtml = '';
        if (input.witness && input.witness.length > 0) {
            const items = input.witness.map((w, wi) => {
                if (!w || w === '') return `<span class="witness-item">[${wi}] (empty)</span>`;
                const display = w.length > 24 ? w.slice(0, 12) + '…' + w.slice(-12) : w;
                return `<span class="witness-item" title="${w}">[${wi}] ${display}</span>`;
            }).join('');
            witnessHtml = `<div class="input-detail witness-detail"><span class="detail-label">Witness:</span> <div class="witness-stack">${items}</div></div>`;
        }

        // Outpoint reference
        const outpointHtml = input.txid ? `<div class="input-detail"><span class="detail-label">Spends:</span> <span class="monospace" title="${input.txid}:${input.vout}">${truncateAddress(input.txid)}:${input.vout}</span></div>` : '';

        div.innerHTML = `
            <div class="flow-item-header">
                <span class="script-badge" title="${getScriptTypeName(input.script_type)}">${input.script_type || 'unknown'}</span>
                <span class="flow-item-value">${formatSats(value)}</span>
            </div>
            <div class="flow-item-address">${truncateAddress(input.address)}</div>
            ${outpointHtml}${seqHtml}${timelockHtml}${witnessHtml}
        `;
        elements.inputsList.appendChild(div);
    });

    // Render outputs
    data.vout.forEach((output, i) => {
        const isOpReturn = output.script_type === 'op_return';
        const div = document.createElement('div');
        div.className = 'flow-item output-item' + (isOpReturn ? ' op-return' : '');

        let content = `
            <div class="flow-item-header">
                <span class="script-badge">${output.script_type || 'unknown'}</span>
                <span class="flow-item-value">${formatSats(output.value_sats)}</span>
            </div>
        `;

        if (isOpReturn) {
            content += `<div class="flow-item-address">Data Output</div>`;
            if (output.op_return_data_utf8) {
                content += `<div class="op-return-data">`;
                if (output.op_return_protocol && output.op_return_protocol !== 'unknown') {
                    content += `<span class="op-return-protocol">[${output.op_return_protocol}]</span> `;
                }
                content += `"${output.op_return_data_utf8}"</div>`;
            } else if (output.op_return_data_hex) {
                content += `<div class="op-return-data">Hex: ${output.op_return_data_hex.slice(0, 40)}${output.op_return_data_hex.length > 40 ? '...' : ''}</div>`;
            }
        } else {
            content += `<div class="flow-item-address">${truncateAddress(output.address)}</div>`;
        }

        div.innerHTML = content;
        elements.outputsList.appendChild(div);
    });

    // Totals
    elements.totalIn.textContent = formatSats(data.total_input_sats);
    elements.totalOut.textContent = formatSats(data.total_output_sats);
    elements.feeValue.textContent = formatSats(data.fee_sats);

    // Balance bar
    const total = data.total_input_sats;
    const outputPct = total > 0 ? (data.total_output_sats / total) * 100 : 0;
    const feePct = total > 0 ? (data.fee_sats / total) * 100 : 0;

    elements.balanceOutputs.style.width = outputPct + '%';
    elements.balanceFee.style.width = feePct + '%';

    // Only show labels if there's enough space
    elements.balanceOutputsLabel.textContent = outputPct > 15 ? formatSats(data.total_output_sats) : '';
    elements.balanceFeeLabel.textContent = feePct > 10 ? formatSats(data.fee_sats) : '';
}

// ===== Tech Section =====
function updateTechSection(data) {
    // Transaction ID
    elements.techTxid.textContent = data.txid;

    // Witness Transaction ID (SegWit only)
    if (data.wtxid) {
        elements.wtxidCard.style.display = 'block';
        elements.techWtxid.textContent = data.wtxid;
    } else {
        elements.wtxidCard.style.display = 'none';
    }

    // Type
    if (data.segwit) {
        elements.techType.innerHTML = '<span style="color: var(--success)">✓ SegWit</span>';
        elements.techTypeHint.textContent = 'Uses newer efficient format';
    } else {
        elements.techType.innerHTML = '<span style="color: var(--text-muted)">Legacy</span>';
        elements.techTypeHint.textContent = 'Original Bitcoin format';
    }

    // Size metrics
    elements.techSize.textContent = data.size_bytes;
    elements.techVbytes.textContent = data.vbytes;
    elements.techWeight.textContent = data.weight;

    // Fee rate
    elements.techFeerate.textContent = data.fee_rate_sat_vb.toFixed(2);
    elements.feeRateHint.textContent = getFeeRateHint(data.fee_rate_sat_vb);

    // Version
    elements.techVersion.textContent = data.version;

    // Locktime
    if (data.locktime_type === 'none' || data.locktime_value === 0) {
        elements.techLocktime.textContent = 'None';
        elements.locktimeHint.textContent = 'No time lock';
    } else if (data.locktime_type === 'block_height') {
        elements.techLocktime.textContent = `Block ${data.locktime_value.toLocaleString()}`;
        elements.locktimeHint.textContent = 'Cannot confirm until this block';
    } else {
        const date = new Date(data.locktime_value * 1000);
        elements.techLocktime.textContent = date.toLocaleDateString();
        elements.locktimeHint.textContent = 'Cannot confirm until this date';
    }

    // RBF
    if (data.rbf_signaling) {
        elements.techRbf.innerHTML = '<span style="color: var(--warning)">⚠ Enabled</span>';
    } else {
        elements.techRbf.innerHTML = '<span style="color: var(--success)">✓ Disabled</span>';
    }

    // SegWit Savings
    if (data.segwit_savings) {
        elements.segwitSavings.style.display = 'block';
        const savings = data.segwit_savings;
        const maxWeight = savings.weight_if_legacy;

        const actualPct = (savings.weight_actual / maxWeight) * 100;
        elements.savingsActual.style.width = actualPct + '%';
        elements.savingsActualLabel.textContent = savings.weight_actual + ' WU';

        elements.savingsLegacy.style.width = '100%';
        elements.savingsLegacyLabel.textContent = savings.weight_if_legacy + ' WU';

        elements.savingsText.innerHTML = `🎉 This transaction saves <strong>${savings.savings_pct.toFixed(1)}%</strong> with SegWit! ` +
            `<span class="savings-breakdown">Witness: ${savings.witness_bytes}B · Non-witness: ${savings.non_witness_bytes}B · Total: ${savings.total_bytes}B</span>`;
    } else {
        elements.segwitSavings.style.display = 'none';
    }
}

// ===== Copy to Clipboard =====
document.querySelectorAll('.copy-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
        const targetId = btn.dataset.copy;
        const targetEl = document.getElementById(targetId);

        if (targetEl) {
            try {
                await navigator.clipboard.writeText(targetEl.textContent);
                btn.textContent = '✓';
                btn.classList.add('copied');
                setTimeout(() => {
                    btn.textContent = '📋';
                    btn.classList.remove('copied');
                }, 1500);
            } catch (err) {
                console.error('Failed to copy:', err);
            }
        }
    });
});

// ===== Tooltips =====
document.querySelectorAll('.term, [data-tooltip]').forEach(el => {
    el.addEventListener('mouseenter', (e) => {
        const tooltip = el.dataset.tooltip || el.getAttribute('data-tooltip');
        if (!tooltip) return;

        elements.tooltip.textContent = tooltip;
        elements.tooltip.classList.add('visible');

        const rect = el.getBoundingClientRect();
        const tooltipRect = elements.tooltip.getBoundingClientRect();

        let left = rect.left + (rect.width / 2) - (tooltipRect.width / 2);
        let top = rect.top - tooltipRect.height - 10;

        // Keep within viewport
        if (left < 10) left = 10;
        if (left + tooltipRect.width > window.innerWidth - 10) {
            left = window.innerWidth - tooltipRect.width - 10;
        }
        if (top < 10) {
            top = rect.bottom + 10;
        }

        elements.tooltip.style.left = left + 'px';
        elements.tooltip.style.top = top + 'px';
    });

    el.addEventListener('mouseleave', () => {
        elements.tooltip.classList.remove('visible');
    });
});

// ===== Keyboard Shortcuts =====
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + Enter to analyze
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        if (elements.inputPanel.style.display !== 'none') {
            elements.btnAnalyze.click();
        }
    }

    // Escape to go back
    if (e.key === 'Escape') {
        if (elements.resultsPanel.style.display !== 'none' ||
            elements.errorPanel.style.display !== 'none') {
            showInputPanel();
        }
    }
});

// ===== Block File Upload =====

// File input display names
['blk', 'rev', 'xor'].forEach(prefix => {
    const input = document.getElementById(`${prefix}-file`);
    const nameEl = document.getElementById(`${prefix}-file-name`);
    if (input && nameEl) {
        input.addEventListener('change', () => {
            nameEl.textContent = input.files.length ? input.files[0].name : `Choose ${prefix}*.dat`;
            if (input.files.length) nameEl.classList.add('file-selected');
            else nameEl.classList.remove('file-selected');
        });
    }
});

// Block back button
const btnBlockBack = document.getElementById('btn-block-back');
if (btnBlockBack) {
    btnBlockBack.addEventListener('click', showInputPanel);
}

// Analyze Block button
const btnAnalyzeBlock = document.getElementById('btn-analyze-block');
if (btnAnalyzeBlock) {
    btnAnalyzeBlock.addEventListener('click', async () => {
        const blkInput = document.getElementById('blk-file');
        const revInput = document.getElementById('rev-file');
        const xorInput = document.getElementById('xor-file');

        if (!blkInput.files.length || !revInput.files.length) {
            showError('MISSING_FILES', 'Please select both a blk*.dat and rev*.dat file.');
            return;
        }

        const progressEl = document.getElementById('block-progress');
        const progressFill = document.getElementById('progress-fill');
        const progressText = document.getElementById('progress-text');

        progressEl.style.display = 'block';
        progressFill.style.width = '20%';
        progressText.textContent = 'Reading files…';

        try {
            const blkData = await readFileAsBase64(blkInput.files[0]);
            progressFill.style.width = '40%';

            const revData = await readFileAsBase64(revInput.files[0]);
            progressFill.style.width = '60%';

            let xorKey = '';
            if (xorInput.files.length) {
                xorKey = await readFileAsBase64(xorInput.files[0]);
            }

            progressFill.style.width = '70%';
            progressText.textContent = 'Analyzing block…';

            const response = await fetch('/api/analyze/block', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    blk_data: blkData,
                    rev_data: revData,
                    xor_key: xorKey
                })
            });

            progressFill.style.width = '90%';

            const result = await response.json();

            progressFill.style.width = '100%';
            progressText.textContent = 'Done!';

            // Hide progress after a brief delay
            setTimeout(() => { progressEl.style.display = 'none'; }, 300);

            if (Array.isArray(result)) {
                displayBlockResults(result);
            } else if (result.ok === false && result.error) {
                showError(result.error.code || 'BLOCK_ERROR', result.error.message || 'Block analysis failed');
            } else {
                // Single block result - wrap in array
                displayBlockResults([result]);
            }
        } catch (err) {
            progressEl.style.display = 'none';
            showError('BLOCK_ERROR', 'Failed to analyze block: ' + err.message);
        }
    });
}

function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
            const arrayBuffer = reader.result;
            const bytes = new Uint8Array(arrayBuffer);
            let binary = '';
            for (let i = 0; i < bytes.byteLength; i++) {
                binary += String.fromCharCode(bytes[i]);
            }
            resolve(btoa(binary));
        };
        reader.onerror = () => reject(new Error('Failed to read file'));
        reader.readAsArrayBuffer(file);
    });
}

function displayBlockResults(blocks) {
    const container = document.getElementById('block-results-content');
    container.innerHTML = '';

    blocks.forEach((block, blockIdx) => {
        const blockDiv = document.createElement('div');
        blockDiv.className = 'block-card';

        if (block.ok === false) {
            blockDiv.innerHTML = `
                <div class="block-header block-error-header">
                    <h2>Block ${blockIdx + 1} — Error</h2>
                    <p class="block-error-msg">${block.error?.message || 'Unknown error'}</p>
                </div>`;
            container.appendChild(blockDiv);
            return;
        }

        const txCount = block.transactions ? block.transactions.length : 0;
        const totalFees = block.transactions
            ? block.transactions.reduce((sum, tx) => sum + (tx.fee_sats || 0), 0)
            : 0;

        // Build script type summary
        const scriptSummary = buildBlockScriptSummary(block.transactions || []);

        blockDiv.innerHTML = `
            <div class="block-header">
                <h2>🧱 Block ${blockIdx + 1}</h2>
                <div class="block-hash monospace">${block.block_hash || 'Unknown'}</div>
            </div>
            <div class="block-overview">
                <div class="block-stat">
                    <span class="block-stat-label">Transactions</span>
                    <span class="block-stat-value">${txCount}</span>
                </div>
                <div class="block-stat">
                    <span class="block-stat-label">Total Fees</span>
                    <span class="block-stat-value">${formatSats(totalFees)}</span>
                </div>
                <div class="block-stat">
                    <span class="block-stat-label">Version</span>
                    <span class="block-stat-value">${block.version || '—'}</span>
                </div>
                <div class="block-stat">
                    <span class="block-stat-label">Timestamp</span>
                    <span class="block-stat-value">${block.timestamp ? new Date(block.timestamp * 1000).toLocaleString() : '—'}</span>
                </div>
                ${block.merkle_root ? `<div class="block-stat block-stat-wide">
                    <span class="block-stat-label">Merkle Root</span>
                    <span class="block-stat-value monospace" style="font-size:0.75rem;">${block.merkle_root}</span>
                </div>` : ''}
            </div>
            ${scriptSummary ? `<div class="block-script-summary">
                <h4>Script Type Summary</h4>
                <div class="script-summary-grid">${scriptSummary}</div>
            </div>` : ''}
            <div class="block-tx-list">
                <h3>Transactions</h3>
                <div class="block-tx-items" id="block-${blockIdx}-txs"></div>
            </div>
        `;

        container.appendChild(blockDiv);

        // Render transaction items
        const txContainer = document.getElementById(`block-${blockIdx}-txs`);
        (block.transactions || []).forEach((tx, txIdx) => {
            const txItem = document.createElement('div');
            txItem.className = 'block-tx-item';

            const isCoinbase = tx.vin && tx.vin.length === 1 && tx.vin[0].txid === '0000000000000000000000000000000000000000000000000000000000000000';

            txItem.innerHTML = `
                <div class="block-tx-header" onclick="this.parentElement.classList.toggle('expanded')">
                    <div class="block-tx-summary">
                        <span class="block-tx-idx">#${txIdx}</span>
                        ${isCoinbase ? '<span class="coinbase-badge">Coinbase</span>' : ''}
                        <span class="block-tx-id monospace">${truncateAddress(tx.txid)}</span>
                    </div>
                    <div class="block-tx-meta">
                        ${tx.segwit ? '<span class="segwit-tag">SegWit</span>' : ''}
                        <span class="block-tx-fee">${tx.fee_sats != null ? formatSats(tx.fee_sats) + ' fee' : ''}</span>
                        <span class="expand-arrow">▸</span>
                    </div>
                </div>
                <div class="block-tx-detail">
                    ${renderBlockTxDetail(tx)}
                </div>
            `;
            txContainer.appendChild(txItem);
        });
    });

    // Show raw JSON toggle
    const rawSection = document.createElement('details');
    rawSection.className = 'raw-json-section';
    rawSection.innerHTML = `<summary>📄 View Raw JSON</summary><pre>${JSON.stringify(blocks, null, 2)}</pre>`;
    container.appendChild(rawSection);

    showBlockResultsPanel();
}

function buildBlockScriptSummary(transactions) {
    const outputTypes = {};
    const inputTypes = {};

    transactions.forEach(tx => {
        (tx.vout || []).forEach(out => {
            const t = out.script_type || 'unknown';
            outputTypes[t] = (outputTypes[t] || 0) + 1;
        });
        (tx.vin || []).forEach(inp => {
            const t = inp.script_type || 'unknown';
            inputTypes[t] = (inputTypes[t] || 0) + 1;
        });
    });

    let html = '';
    const allTypes = new Set([...Object.keys(outputTypes), ...Object.keys(inputTypes)]);
    if (allTypes.size === 0) return '';

    allTypes.forEach(type => {
        const outCount = outputTypes[type] || 0;
        const inCount = inputTypes[type] || 0;
        html += `<div class="script-summary-item">
            <span class="script-badge">${type}</span>
            <span class="script-summary-counts">${inCount} in / ${outCount} out</span>
        </div>`;
    });
    return html;
}

function renderBlockTxDetail(tx) {
    let html = '<div class="block-tx-grid">';

    // Inputs
    html += '<div class="block-tx-inputs"><h5>Inputs</h5>';
    (tx.vin || []).forEach((inp, i) => {
        const val = inp.prevout?.value_sats;
        html += `<div class="mini-flow-item">
            <span class="script-badge">${inp.script_type || '?'}</span>
            <span class="mini-addr">${truncateAddress(inp.address)}</span>
            ${val != null ? `<span class="mini-val">${formatSats(val)}</span>` : ''}
        </div>`;
    });
    html += '</div>';

    // Outputs
    html += '<div class="block-tx-outputs"><h5>Outputs</h5>';
    (tx.vout || []).forEach((out, i) => {
        const isOp = out.script_type === 'op_return';
        html += `<div class="mini-flow-item ${isOp ? 'op-return' : ''}">
            <span class="script-badge">${out.script_type || '?'}</span>
            <span class="mini-addr">${isOp ? 'OP_RETURN' : truncateAddress(out.address)}</span>
            <span class="mini-val">${formatSats(out.value_sats)}</span>
        </div>`;
    });
    html += '</div>';

    html += '</div>';

    // Metadata line
    html += `<div class="block-tx-meta-detail">
        Size: ${tx.size_bytes || '?'}B | Weight: ${tx.weight || '?'} WU | vBytes: ${tx.vbytes || '?'} | Fee rate: ${tx.fee_rate_sat_vb != null ? tx.fee_rate_sat_vb.toFixed(2) : '?'} sat/vB
    </div>`;

    return html;
}

// ===== Initialize =====
console.log('⛓️ Chain Lens initialized');
