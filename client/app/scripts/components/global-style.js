import { createGlobalStyle } from 'styled-components';
import { color } from 'weaveworks-ui-components/lib/theme/selectors';

const GlobalStyle = createGlobalStyle`
  div {
    background: ${props => props.theme.scope.background};
  }

  /**
  * Copyright (c) 2014 The xterm.js authors. All rights reserved.
  * Copyright (c) 2012-2013, Christopher Jeffrey (MIT License)
  * https://github.com/chjj/term.js
  * @license MIT
  *
  * Permission is hereby granted, free of charge, to any person obtaining a copy
  * of this software and associated documentation files (the "Software"), to deal
  * in the Software without restriction, including without limitation the rights
  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  * copies of the Software, and to permit persons to whom the Software is
  * furnished to do so, subject to the following conditions:
  *
  * The above copyright notice and this permission notice shall be included in
  * all copies or substantial portions of the Software.
  *
  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  * THE SOFTWARE.
  *
  * Originally forked from (with the author's permission):
  *   Fabrice Bellard's javascript vt100 for jslinux:
  *   http://bellard.org/jslinux/
  *   Copyright (c) 2011 Fabrice Bellard
  *   The original design remains. The terminal itself
  *   has been extended to include xterm CSI codes, among
  *   other features.
  */

  /**
  *  Default styles for xterm.js
  */

  .xterm {
      font-family: $font-family-monospace;
      font-feature-settings: "liga" 0;
      position: relative;
      user-select: none;
      /* stylelint-disable property-no-vendor-prefix */
      -ms-user-select: none;
      -webkit-user-select: none;
      /* stylelint-enable property-no-vendor-prefix */
  }

  .xterm.focus,
  .xterm:focus {
      outline: none;
  }

  .xterm .xterm-helpers {
      position: absolute;
      top: 0;
      /**
      * The z-index of the helpers must be higher than the canvases in order for
      * IMEs to appear on top.
      */
      /* stylelint-disable sh-waqar/declaration-use-variable */
      z-index: 10;
      /* stylelint-enable sh-waqar/declaration-use-variable */
  }

  .xterm .xterm-helper-textarea {
      /*
      * HACK: to fix IE's blinking cursor
      * Move textarea out of the screen to the far left, so that the cursor is not visible.
      */
      position: absolute;
      opacity: 0;
      left: -9999em;
      top: 0;
      width: 0;
      height: 0;
      /* stylelint-disable sh-waqar/declaration-use-variable */
      z-index: -10;
      /* stylelint-enable sh-waqar/declaration-use-variable */
      /** Prevent wrapping so the IME appears against the textarea at the correct position */
      white-space: nowrap;
      overflow: hidden;
      resize: none;
  }

  .xterm .composition-view {
      /* TODO: Composition position got messed up somewhere */
      background: ${color('black')};
      color: ${color('white')};
      display: none;
      position: absolute;
      white-space: nowrap;
      z-index: ${props => props.theme.layers.front};
  }

  .xterm .composition-view.active {
      display: block;
  }

  .xterm .xterm-viewport {
      /* On OS X this is required in order for the scroll bar to appear fully opaque */
      background-color: ${color('black')};
      overflow-y: scroll;
      cursor: default;
      position: absolute;
      right: 0;
      left: 0;
      top: 0;
      bottom: 0;
  }

  .xterm .xterm-screen {
      position: relative;
  }

  .xterm .xterm-screen canvas {
      position: absolute;
      left: 0;
      top: 0;
  }

  .xterm .xterm-scroll-area {
      visibility: hidden;
  }

  .xterm-char-measure-element {
      display: inline-block;
      visibility: hidden;
      position: absolute;
      top: 0;
      left: -9999em;
      line-height: normal;
  }

  .xterm.enable-mouse-events {
      /* When mouse events are enabled (eg. tmux), revert to the standard pointer cursor */
      cursor: default;
  }

  .xterm:not(.enable-mouse-events) {
      cursor: text;
  }

  .xterm .xterm-accessibility,
  .xterm .xterm-message {
      position: absolute;
      left: 0;
      top: 0;
      bottom: 0;
      right: 0;
      /* stylelint-disable sh-waqar/declaration-use-variable */
      z-index: 100;
      /* stylelint-enable sh-waqar/declaration-use-variable */
      color: transparent;
  }

  .xterm .live-region {
      position: absolute;
      left: -9999px;
      width: 1px;
      height: 1px;
      overflow: hidden;
  }

  .xterm-cursor-pointer {
      cursor: pointer;
  }
`;

export default GlobalStyle;
