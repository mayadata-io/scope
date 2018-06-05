import React from 'react';
import { connect } from 'react-redux';
import classNames from 'classnames';

import { enterEdge, leaveEdge } from '../actions/app-actions';
import { encodeIdAttribute, decodeIdAttribute } from '../utils/dom-utils';

function IsStorageComponent(id) {
  if (id === 'persietent_volume' || id === 'storage_class' || id === 'persistent_volume_claim') {
    return true;
  }
  return false;
}

// getAdjacencyClass takes id which contains information about edge as a topology
// of parent and child node.

function getAdjacencyClass(id) {
  const topologyId = id.split('---');
  const from = topologyId[0].split(';');
  const to = topologyId[1].split(';');
  if (from[1] !== undefined && to[1] !== undefined) {
    from[1] = from[1].slice(1, -1);
    to[1] = to[1].slice(1, -1);
    if (IsStorageComponent(from[1]) || IsStorageComponent(to[1])) {
      return 'link-storage';
    }
  }
  return 'link-none';
}

class Edge extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleMouseEnter = this.handleMouseEnter.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  render() {
    const {
      id, path, highlighted, focused, thickness, source, target
    } = this.props;
    const shouldRenderMarker = (focused || highlighted) && (source !== target);
    const className = classNames('edge', { highlighted });
    return (
      <g
        id={encodeIdAttribute(id)}
        className={className}
        onMouseEnter={this.handleMouseEnter}
        onMouseLeave={this.handleMouseLeave}
      >
        <path className="shadow" d={path} style={{ strokeWidth: 10 * thickness }} />
        <path
          className={getAdjacencyClass(id)}
          d={path}
          style={{ strokeWidth: 5 }}
        />
        <path
          className="link"
          d={path}
          markerEnd={shouldRenderMarker ? 'url(#end-arrow)' : null}
          style={{ strokeWidth: thickness }}
        />
      </g>
    );
  }

  handleMouseEnter(ev) {
    this.props.enterEdge(decodeIdAttribute(ev.currentTarget.id));
  }

  handleMouseLeave(ev) {
    this.props.leaveEdge(decodeIdAttribute(ev.currentTarget.id));
  }
}

function mapStateToProps(state) {
  return {
    contrastMode: state.get('contrastMode')
  };
}

export default connect(
  mapStateToProps,
  { enterEdge, leaveEdge }
)(Edge);
