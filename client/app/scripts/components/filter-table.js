import React from 'react';
import { connect } from 'react-redux';
import Checkbox from './checkbox';
import { canvasMarginsSelector } from '../selectors/canvas';
import { hideFilter } from '../actions/app-actions';


const OPTIONS = ['Access modes', 'Iqn', 'Logical Sector Size',
  'Lowest Temperature', 'Model', 'Percent Endurance Used', 'Physical Sector Size', 'Provisioner', 'Replication Factor', 'Rotation Rate',
  'Serial', 'Storage class', 'Storage driver', 'Total Bytes Written', 'Type', 'Vendor', 'Volume',
  'Iops(R)', 'Iops(W)', 'Latency(R)', 'Latency(W)', 'Throughput(R)', 'Throughput(W)'];

export const tmp = ['docker_container_ports', 'Current Temperature', 'Device Utilization Rate', 'docker_container_id', 'docker_image_id',
  'docker_container_command', 'docker_container_networks', 'docker_container_networks', 'Firmware Revision', 'Memory', 'Load (1m)', 'CPU', 'Highest Temperature', 'Capacity', 'State', 'Volume claim', 'Status',
  '# Threads', 'Command', 'PID', 'Parent PID', 'Created', 'IPs', 'Image name', 'Image tag', 'Restart #', 'Uptime', 'IP', 'Namespace', 'Observed gen.',
];


class FilterTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      checkboxes: OPTIONS.reduce(
        (options, option) => ({
          ...options,
          [option]: false
        }),
        {}
      )
    };
  }

  selectAllCheckboxes = (isSelected) => {
    Object.keys(this.state.checkboxes).forEach((checkbox) => {
      this.setState(prevState => ({
        checkboxes: {
          ...prevState.checkboxes,
          [checkbox]: isSelected
        }
      }));
    });
  };

  selectAll = () => this.selectAllCheckboxes(true);

  selectCheckboxes = (isSelected) => {
    Object.keys(this.state.checkboxes).forEach((checkbox) => {
      this.setState(checkbox = {
        ...checkbox,
        [checkbox]: isSelected
      });
      tmp.pop(checkbox);
    });
  };

  deselectAll = () => this.selectCheckboxes(false);

  handleCheckboxChange = (changeEvent) => {
    const { name } = changeEvent.target;

    this.setState(prevState => ({
      checkboxes: {
        ...prevState.checkboxes,
        [name]: !prevState.checkboxes[name]
      }
    }));
  };

  handleFormSubmit = (formSubmitEvent) => {
    formSubmitEvent.preventDefault();

    Object.keys(this.state.checkboxes)
      .filter(checkbox => this.state.checkboxes[checkbox])
      .forEach((checkbox) => {
        tmp.push(checkbox);
      });
  };

  createCheckbox = option => (
    <Checkbox
      label={option}
      isSelected={this.state.checkboxes[option]}
      onCheckboxchecked={this.handleCheckboxChange}
      key={option}
    />
  );
  createCheckboxes = () => OPTIONS.map(this.createCheckbox);

  render() {
    const { canvasMargins, onClickClose } = this.props;
    return (
      <div className="help-panel-wrapper">
        <div className="help-panel" style={{marginTop: canvasMargins.top}}>
          <div className="help-panel-header">
            <h2>Filter</h2>
          </div>
          <div className="help-panel-main">
            <div className="help-panel-fields">
              <h2>Filter Columns </h2>
              <p>
       Filter columns in the currently <br />
       selected {} topology:
              </p>
              <div className="help-panel-fields-fields">
                <div className="help-panel-fields-fields-column">
                  <h3>Columns</h3>

                  <form onSubmit={this.handleFormSubmit}>
                    <div className="help-panel-fields-fields-column-content">
                      {this.createCheckboxes()}
                    </div>
                    <div className="form-group mt-2">
                      <button
                        type="button"
                        className="tour-step-anchor view-mode-selector-action view-Resources-action"
                        onClick={this.selectAll}
                >
                  Select All
                      </button>
                      <button
                        type="button"
                        className="tour-step-anchor view-mode-selector-action view-Resources-action"
                        onClick={this.deselectAll}
                >
                  Reset
                      </button>
                      <button type="submit" className="tour-step-anchor view-mode-selector-action view-Resources-action">
                  Show
                      </button>

                    </div>
                  </form>
                  <br />
                  <div className="help-panel-tools">
                    <i
                      title="Close details"
                      className="fa fa-times"
                      onClick={onClickClose}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
function mapStateToProps(state) {
  return {
    canvasMargins: canvasMarginsSelector(state)
  };
}


export default connect(mapStateToProps, {
  onClickClose: hideFilter
})(FilterTable);
