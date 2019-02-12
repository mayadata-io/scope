import React from 'react';

const checkbox = ({ label, isSelected, onCheckboxchecked }) => (
  <div className="form-check">
    <input
      type="checkbox"
      name={label}
      checked={isSelected}
      onChange={onCheckboxchecked}
      className="form-check-input"
      />
    {label}
  </div>
);

export default checkbox;
