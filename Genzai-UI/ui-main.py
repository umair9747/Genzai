import streamlit as st
import requests
import json
import time
from datetime import datetime
import plotly.express as px
import pandas as pd
from pathlib import Path

# Define the dark neon green theme colors
THEME = {
    'primary': '#00FF9D',      # Neon green
    'secondary': '#00CC7E',    # Darker neon green
    'background': '#0A1A14',   # Very dark green-tinted black
    'surface': '#132921',      # Dark green surface
    'card': '#1A3329',         # Slightly lighter dark green
    'text': '#E0FFF1',         # Light green-tinted white
    'subtext': '#7EA894',      # Muted green
    'accent': '#00FFB3',       # Bright neon green accent
    'error': '#FF4B4B',        # Red for errors/alerts
    'warning': '#FFB300',      # Amber for warnings
}

# Category Colors for the visualization
CATEGORY_COLORS = {
    "BuildingAccessControlSystem": "#4CAF50",
    "Camera": "#2196F3",
    "ClimateControl": "#FF9800",
    "HMI": "#9C27B0",
    "HomeAutomation": "#00BCD4",
    "IndustrialAutomation": "#F44336",
    "NAS": "#607D8B",
    "Router": "#3F51B5",
    "SmartPowerControl": "#FFEB3B",
    "VideoSurveillanceandSecurity": "#E91E63",
    "WaterTreatmentSystem": "#009688"
}

# Category Icons (simplified)
CATEGORY_ICONS = {
    "BuildingAccessControlSystem": "🏢",
    "Camera": "📹",
    "ClimateControl": "🌡️",
    "HMI": "🖥️",
    "HomeAutomation": "🏠",
    "IndustrialAutomation": "🏭",
    "NAS": "💾",
    "Router": "📡",
    "SmartPowerControl": "⚡",
    "VideoSurveillanceandSecurity": "👁️",
    "WaterTreatmentSystem": "💧"
}


# Set page configuration
st.set_page_config(
    page_title="Genzai - IoT Security Toolkit",
    page_icon="👁️",
    layout="wide",
    initial_sidebar_state="collapsed"
)

# Enhanced Custom CSS with neon theme and bigger expander titles
st.markdown(f"""
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Chakra+Petch:wght@400;600;700&display=swap');
        
        /* Main theme styles */
        .stApp {{
            background-color: {THEME['background']};
            color: {THEME['text']};
            font-family: 'Chakra Petch', sans-serif;
        }}
        
        /* Header styles */
        .header-container {{
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 20px;
            padding: 1rem;
        }}
        
        .main-title {{
            font-family: 'Chakra Petch', sans-serif;
            font-size: 3.5rem;
            font-weight: 700;
            color: {THEME['primary']};
            text-shadow: 0 0 10px {THEME['primary']};
            text-align: center
        }}
        
        /* Card styles */
        .custom-card {{
            background-color: {THEME['card']};
            border-radius: 15px;
            padding: 2rem;
            margin: 1.5rem 0;
            border: 1px solid {THEME['secondary']};
            box-shadow: 0 0 15px rgba(0, 255, 157, 0.1);
        }}
        
        /* Metric card styles */
        .metric-card {{
            background-color: {THEME['surface']};
            border-radius: 12px;
            padding: 1.5rem;
            text-align: center;
            border: 1px solid {THEME['secondary']};
        }}
        
        .metric-value {{
            font-size: 2.5rem;
            font-weight: 700;
            margin: 0.5rem 0;
            color: {THEME['primary']};
        }}
        
        .metric-label {{
            font-size: 1rem;
            color: {THEME['subtext']};
            font-weight: 600;
        }}
        
        /* Enhanced expander styles */
        .streamlit-expanderHeader  {{
            font-size: 24px !important;  /* Increased from default */
            font-weight: 700 !important;
            background-color: {THEME['surface']} !important;
            border-radius: 10px !important;
            padding: 1rem !important;
            margin-bottom: 1rem !important;
            color: {THEME['primary']} !important;
        }}
        
        /* Status badge styles */
        .status-badge {{
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 0.9rem;
            font-weight: 600;
            display: inline-block;
            margin: 0.5rem 0;
        }}
        
        .status-critical {{
            background-color: rgba(255, 75, 75, 0.15);
            color: #FF4B4B;
            border: 1px solid #FF4B4B;
        }}
        
        .status-warning {{
            background-color: rgba(255, 179, 0, 0.15);
            color: #FFB300;
            border: 1px solid #FFB300;
        }}
    </style>
""", unsafe_allow_html=True)


# Initialize session state if not exists
if 'scan_results' not in st.session_state:
    st.session_state.scan_results = None
    

# Header with logo and title
# logo_path = Path("logo.png")
header_col1, header_col2, header_col3 = st.columns([1,2,1])
with header_col2:
    st.markdown('<h1 class="main-title">  -- GENZAI --</h1>', unsafe_allow_html=True)
    st.markdown('<p style="text-align: center; font-size: 1.3rem; color: var(--subtext);">IoT Security Toolkit</p>', unsafe_allow_html=True)

# Main content
st.markdown("-----")
st.markdown(f'<h3 style="color: {THEME["primary"]}; font-size: 1.8rem;">🎯 Target Input</h3>', unsafe_allow_html=True)

# Create two columns for input methods
col1, col2 = st.columns(2)

st.markdown(f'<p style="color: {THEME["text"]};">Enter target URLs manually:</p>', unsafe_allow_html=True)
targets_input = st.text_area(
    "Enter target URLs (one per line)", 
    height=100,
    help="Enter one target URL per line",
    label_visibility="collapsed"
)

with st.expander("Genzai, upload input file or change API endpoint:"):
    st.markdown(f'<p style="color: {THEME["text"]};">Upload a text file:</p>', unsafe_allow_html=True)
    uploaded_file = st.file_uploader("Choose a text file", type=['txt'], label_visibility="collapsed")
    if uploaded_file is not None:
        targets_input = uploaded_file.getvalue().decode("utf-8")

    api_endpoint = st.text_input("API Endpoint", value="http://0.0.0.0:8080/scan")

if st.button("Start Scan 🚀"):
    if targets_input:
        targets = [target.strip() for target in targets_input.split('\n') if target.strip()]
        targets = list(set(targets))
        
        progress_bar = st.progress(0)
        status = st.empty()
        
        start_time = time.time()
        
        try:
            status.markdown('<div class="custom-card">🔍 Scanning targets...</div>', unsafe_allow_html=True)
            response = requests.post(
                api_endpoint,
                headers={"Content-Type": "application/json"},
                json={"targets": targets}
            )
            
            execution_time = time.time() - start_time
            
            if response.status_code == 200:
                results = response.json()
                # Store results in session state
                st.session_state.scan_results = {
                    'results': results,
                    'targets': targets,
                    'execution_time': execution_time
                }
                status.empty()
                progress_bar.progress(100)
            else:
                st.error(f"Error: API returned status code {response.status_code}")
                
        except requests.exceptions.RequestException as e:
            st.error(f"Error connecting to API: {str(e)}")
    else:
        st.warning("Please enter at least one target URL")

# Display results if they exist in session state
if st.session_state.scan_results:
    results = st.session_state.scan_results['results']
    targets = st.session_state.scan_results['targets']
    execution_time = st.session_state.scan_results['execution_time']
    
    # Metrics Dashboard
    st.markdown("-----")
    col1, col2, col3, col4 = st.columns(4)
    
    with col1:
        st.markdown(f"""
        <div class="metric-card">
            <div class="metric-label">Total Targets</div>
            <div class="metric-value">🎯 {len(targets)}</div>
        </div>
        """, unsafe_allow_html=True)
    
    with col2:
        st.markdown(f"""
        <div class="metric-card">
            <div class="metric-label">Execution Time</div>
            <div class="metric-value">⚡ {execution_time:.2f}s</div>
        </div>
        """, unsafe_allow_html=True)
    
    with col3:
        total_issues = sum(len(result["Issues"]) for result in results["Results"])
        st.markdown(f"""
        <div class="metric-card">
            <div class="metric-label">Issues Found</div>
            <div class="metric-value">🚨 {total_issues}</div>
        </div>
        """, unsafe_allow_html=True)
    
    with col4:
        unique_categories = len(set(result["category"] for result in results["Results"]))
        st.markdown(f"""
        <div class="metric-card">
            <div class="metric-label">Device Categories</div>
            <div class="metric-value">📊 {unique_categories}</div>
        </div>
        """, unsafe_allow_html=True)
    
    # Device Category Distribution visualization
    st.markdown("-----")
    st.subheader("📊 Device Category Distribution")
    
    categories = {}
    for result in results["Results"]:
        cat = result["category"]
        categories[cat] = categories.get(cat, 0) + 1
    
    fig = px.sunburst(
        names=list(categories.keys()),
        parents=[""] * len(categories),
        values=list(categories.values()),
        color=list(categories.keys()),
        color_discrete_map=CATEGORY_COLORS
    )
    
    fig.update_layout(
        paper_bgcolor='rgba(0,0,0,0)',
        plot_bgcolor='rgba(0,0,0,0)',
        font=dict(color='white'),
        margin=dict(t=0, l=0, r=0, b=0)
    )
    
    st.plotly_chart(fig, use_container_width=True)
    
    # Detailed Results
    st.markdown("-----")
    st.subheader("🔍 Detailed Findings")
    
    # Show only first result expanded
    for idx, result in enumerate(results["Results"]):
        with st.expander(
            f"{CATEGORY_ICONS.get(result['category'], '🔍')} {result['Target']}", 
            expanded=(idx == 0)  # Only expand first result
        ):
            st.markdown(f"<h3 style='color: {THEME['primary']}; font-size: 1.8rem; margin: 2rem 0;'>🔍 Device Information</h3>", unsafe_allow_html=True)
            st.markdown(f"""
                <div style='background-color: {THEME["surface"]}; padding: 2rem; border-radius: 12px; margin-top: 1rem;'>
                    <p style='font-size: 1.2rem; color: {THEME["text"]};'><strong>Type:</strong> {result['IoTidentified']}</p>
                    <p style='font-size: 1.2rem; color: {THEME["text"]};'><strong>Category:</strong> {result['category']}</p>
                </div>
            """, unsafe_allow_html=True)
            
            if result["Issues"]:
                st.markdown(f"<h3 style='color: {THEME['primary']}; font-size: 1.8rem; margin: 2rem 0;'>🚨 Security Issues</h3>", unsafe_allow_html=True)
                for issue in result["Issues"]:
                    severity = "critical" if "Authentication" in issue["IssueTitle"] or "Default Password" in issue["IssueTitle"] else "warning"
                    st.markdown(f"""
                        <div style='background-color: {THEME["surface"]}; padding: 1.5rem; border-radius: 12px; margin: 1rem 0;'>
                            <span class='status-badge status-{severity}'>
                                {severity.upper()}
                            </span>
                            <h4 style='color: {THEME["primary"]}; font-size: 1.4rem; margin: 1rem 0;'>{issue['IssueTitle']}</h4>
                            <p style='color: {THEME["text"]}; margin: 0.5rem 0;'><strong>URL:</strong> <code style='background-color: {THEME["background"]}; padding: 0.3rem 0.6rem; border-radius: 4px;'>{issue['URL']}</code></p>
                            <p style='color: {THEME["text"]}; margin: 0.5rem 0;'><strong>Context:</strong> {issue['AdditionalContext']}</p>
                        </div>
                    """, unsafe_allow_html=True)
    
    # Export options
    st.markdown("-----")
    col1, col2 = st.columns([1, 1])
    
    with col1:
        st.download_button(
            label="📥 Export Results (JSON)",
            data=json.dumps(results, indent=2),
            file_name=f"genzai_results_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json",
            mime="application/json"
        )
    
    with col2:
        # Create CSV export
        csv_data = []
        for result in results["Results"]:
            for issue in result["Issues"]:
                csv_data.append({
                    "Target": result["Target"],
                    "Device": result["IoTidentified"],
                    "Category": result["category"],
                    "Issue": issue["IssueTitle"],
                    "URL": issue["URL"]
                })
        
        csv_df = pd.DataFrame(csv_data)
        st.download_button(
            label="📊 Export Results (CSV)",
            data=csv_df.to_csv(index=False),
            file_name=f"genzai_results_{datetime.now().strftime('%Y%m%d_%H%M%S')}.csv",
            mime="text/csv"
        )

## Footer
st.markdown(f"""
    <div style='text-align: center; padding: 2rem; margin-top: 3rem; border-top: 1px solid {THEME["secondary"]};'>
        <p style='color: {THEME["subtext"]}; font-size: 1.1rem;'>Live to you by: rumble773 & umair9747</p>
        <p style='color: {THEME["subtext"]}; font-size: 0.9rem;'>Genzai v2.0.0</p>
    </div>
""", unsafe_allow_html=True)