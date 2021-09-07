
function UnitBezier(p1x, p1y, p2x, p2y) {
    // Calculate the polynomial coefficients, implicit first and last control points are (0,0) and (1,1).
    this.cx = 3.0 * p1x;
    this.bx = 3.0 * (p2x - p1x) - this.cx;
    this.ax = 1.0 - this.cx - this.bx;

    this.cy = 3.0 * p1y;
    this.by = 3.0 * (p2y - p1y) - this.cy;
    this.ay = 1.0 - this.cy - this.by;

    this.p1x = p1x;
    this.p1y = p1y;
    this.p2x = p2x;
    this.p2y = p2y;
}

UnitBezier.prototype = {
    sampleCurveX: function (t) {
        // `ax t^3 + bx t^2 + cx t' expanded using Horner's rule.
        return ((this.ax * t + this.bx) * t + this.cx) * t;
    },

    sampleCurveY: function (t) {
        return ((this.ay * t + this.by) * t + this.cy) * t;
    },

    sampleCurveDerivativeX: function (t) {
        return (3.0 * this.ax * t + 2.0 * this.bx) * t + this.cx;
    },

    solveCurveX: function (x, epsilon) {
        if (epsilon === undefined) epsilon = 1e-6;

        if (x < 0.0) return 0.0;
        if (x > 1.0) return 1.0;

        var t = x;

        // First try a few iterations of Newton's method - normally very fast.
        for (var i = 0; i < 8; i++) {
            var x2 = this.sampleCurveX(t) - x;
            if (Math.abs(x2) < epsilon) return t;

            var d2 = this.sampleCurveDerivativeX(t);
            if (Math.abs(d2) < 1e-6) break;

            t = t - x2 / d2;
        }

        // Fall back to the bisection method for reliability.
        var t0 = 0.0;
        var t1 = 1.0;
        t = x;

        for (i = 0; i < 20; i++) {
            x2 = this.sampleCurveX(t);
            if (Math.abs(x2 - x) < epsilon) break;

            if (x > x2) {
                t0 = t;
            } else {
                t1 = t;
            }

            t = (t1 - t0) * 0.5 + t0;
        }

        return t;
    },

    solve: function (x, epsilon) {
        return this.sampleCurveY(this.solveCurveX(x, epsilon));
    }
};

class World {
    constructor(canvas) {
        this.lat = 0
        this.lng = 0
        this.zoom = 0
        let dpr = window.devicePixelRatio
        canvas.width = canvas.parentNode.clientWidth * dpr
        canvas.height = canvas.parentNode.clientHeight * dpr
        // canvas.width = 3840
        // canvas.height = 2160
        let ctx = canvas.getContext('2d',{alpha:false})
        this.context = ctx
        this.canvas = canvas
        ctx.imageSmoothingEnabled = false
        // this.context.scale(dpr,dpr)
        ctx.fillStyle = "ghostwhite"
        ctx.fillRect(0,0,canvas.width,canvas.height)

        let pattern = document.createElement("canvas")
        pattern.width = 16
        pattern.height = 16
        let pctx = pattern.getContext("2d")
        pctx.imageSmoothingEnabled = false
        pctx.strokeStyle = "black"
        pctx.beginPath()
        pctx.moveTo(0,0)
        pctx.lineTo(0,16)
        pctx.lineTo(16,16)
        pctx.stroke()

        this.pattern = ctx.createPattern(pattern,'repeat')


        let tile = document.createElement("canvas")
        tile.width = 1024 * dpr
        tile.height = 1024 * dpr
        let tileCtx = tile.getContext('2d')
        tileCtx.scale(dpr,dpr)
        tileCtx.fillStyle = "steelblue"
        tileCtx.fillRect(0,0,1024,1024)
        this.tile = tile

        let label =document.createElement("canvas")
        label.width = 200
        label.height = 40
        let labelCtx = label.getContext('2d')
        labelCtx.scale(dpr,dpr)
        labelCtx.strokeStyle = "red"
        labelCtx.lineWidth = 2
        labelCtx.font = "600 12px sans-serif"
        labelCtx.fillRect(0,0,100,20)
        labelCtx.strokeText("foobarfoobarfoobarfoobarfoobar",0,16)
        labelCtx.fillText("foobarfoobarfoobarfoobarfoobar",0,16)
        this.label = label

        this.positions = []
        for (var x = 0; x < 20; x++) {
            for (var y = 0; y < 20; y++) {
                this.positions.push([x*40,y*40])

            }
        }

        this.needsUpdate = true
    } 

    project = (in_lat,in_lng) => {
        var d = Math.PI / 180,
            max = 85.0511287798,
            lat = Math.max(Math.min(max, in_lat), -max),
            sin = Math.sin(lat * d);

        return [in_lng * d,Math.log((1 + sin) / (1 - sin)) / 2]
    }


    resize = () => {
        let dpr = window.devicePixelRatio
        this.canvas.width = this.canvas.parentNode.clientWidth * dpr
        this.canvas.height = this.canvas.parentNode.clientHeight * dpr
    }

    draw(begin,end) {
        // at zoom 0, the 6378137*2 diameter is 256 pixels
        this.context.clearRect(0,0,this.canvas.width,this.canvas.height)
        let dpr = window.devicePixelRatio
        begin()
        let cx = this.canvas.width / 2 / dpr
        let cy = this.canvas.height / 2 / dpr
        var factor = Math.pow(2,this.zoom)
        let projected = this.project(this.lat,this.lng)
        let tx = projected[0] / Math.PI / 2 * 256 * factor
        let ty = projected[1] / Math.PI / 2 * 256 * factor

        let sz = 256*factor
        this.context.save()
        this.context.translate(-tx+cx,ty+cy)
        this.context.fillStyle = "red"
        this.context.fillRect(-sz/2,-sz/2,sz,sz)
        this.context.fillStyle = this.pattern
        this.context.scale(factor,factor)
        this.context.fillRect(-sz/2/factor,-sz/2/factor,sz/factor,sz/factor)
        this.context.restore()

        // this.context.translate(this.dx,this.dy)
        // // draw all tiles
        // this.context.drawImage(this.tile,0,0,1024* this.zoom,1024 * this.zoom)

        // for (var position of this.positions) {
        //     this.context.drawImage(this.label,position[0] * this.zoom,position[1] * this.zoom)
        //     // console.log(position)
        // }
        end()

        // draw all labels
    }
}
