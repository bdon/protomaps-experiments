// https://github.com/mariusandra/pigeon-maps
// # The MIT License (MIT)
// Copyright (c) 2016 Marius Andra <marius.andra@gmail.com>
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

const ANIMATION_TIME = 300;
const SCROLL_PIXELS_FOR_ZOOM_LEVEL = 150;

let easeOutQuad = (t) => {
  return t * (2 - t);
};

class Scroller {
    constructor() {
        this._isAnimating = false;
        this._lastCenter = {lat:0,lng:0};
        this._lastZoom = 0;
    }

    wheel(event) {
        const addToZoom = -event.deltaY / SCROLL_PIXELS_FOR_ZOOM_LEVEL;
        this.zoomAroundMouse(addToZoom,event);
    }

    zoomAroundMouse(zoomDiff,event) {
        let zoom = this.map.getZoom()
        let zoomTarget = zoom + zoomDiff; 
        let latLngNow = this.map.getCenter()
        zoomTarget = zoomDiff < 0 ? Math.floor(zoomTarget) : Math.ceil(zoomTarget);
        zoomTarget = Math.max(this.map.getMinZoom(), Math.min(zoomTarget, this.map.getMaxZoom()));
        this.setCenterZoomTarget(zoomTarget, latLngNow)
    }

    setCenterZoomTarget(zoom,zoomAround) {
        if (this._isAnimating) {
            cancelAnimationFrame(this._animFrame)
            const { centerStep, zoomStep } = this.animationStep(performance.now())
            this._centerStart = centerStep
            this._zoomStart = zoomStep
        } else {
            this._isAnimating = true
            this._zoomStart = this._lastZoom
        }

        this._animationStart = performance.now()
        this._animationEnd = this._animationStart + 300
        this._zoomTarget = zoom
        this._animFrame = requestAnimationFrame(this.animate)
    }

    animate = (timestamp) => {
        if (!this._animationEnd || timestamp >= this._animationEnd) {
            this._isAnimating = false
            this.map._moveEnd(true);

        } else {
            const { centerStep, zoomStep } = this.animationStep(timestamp);
            this.setCenterZoom(centerStep, zoomStep)
            this._animFrame = requestAnimationFrame(this.animate)
        }
    }

    setCenterZoom(center,zoom) {
        this.map._move(center, zoom);
        this._lastZoom = zoom;
        this._lastCenter = center;
    }

    animationStep = (timestamp) => {
        const length = this._animationEnd - this._animationStart;
        const progress = Math.max(timestamp - this._animationStart, 0);
        const percentage = easeOutQuad(progress / length);
        const zoomDiff = (this._zoomTarget - this._zoomStart) * percentage;
        const zoomStep = this._zoomStart + zoomDiff;
        return { centerStep: {lat:0,lng:0}, zoomStep:zoomStep }
    }
}

const leafletHandler = () => {
    let s = new Scroller();

    return L.Handler.extend({
        addHooks: function() {
            s.map = this._map

            L.DomEvent.on(this._map._container, 'wheel', this._wheel, this);
        },

        removeHooks: function() {
            L.DomEvent.off(this._map._container, 'wheel', this._wheel, this);
        },

        _wheel: function(event) {
            s.wheel(event);
        }
    })
}