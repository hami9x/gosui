var canvas = new fabric.Canvas("canvas");
canvas.selection = false;
canvas.forEachObject(function(o) {
  o.selectable = false;
});

function gosuiAddUnsel(canvas, e) {
	canvas.add(e)
	e.selectable = false
}

function gosuiSameRadiisArray(rad) {
    var a = [];
    for (var i=0; i<4; i++) {
      a[i]=rad;
    }
    return a;
}


function gosuiCanvasRRect(ctx, x, y, w, h, radiis){
  var tl = radiis[0], tr = radiis[1], br = radiis[2], bl = radiis[3];
  var r = x + w,
      b = y + h;

  ctx.beginPath();
  ctx.moveTo(x+tl, y);
  ctx.lineTo(r-(tr), y);
  ctx.quadraticCurveTo(r, y, r, y+tr);
  ctx.lineTo(r, b-br);
  ctx.quadraticCurveTo(r, b, r-br, b);
  ctx.lineTo(x+bl, b);
  ctx.quadraticCurveTo(x, b, x, b-bl);
  ctx.lineTo(x, y+tl);
  ctx.quadraticCurveTo(x, y, x+tl, y);
}

var RoundedRect = fabric.util.createClass(fabric.Rect, {

  type: 'roundedRect',

  initialize: function(options) {
    options || (options = { });

    this.callSuper('initialize', options);
    this.set('cornerRadiis', options.cornerRadiis || gosuiSameRadiisArray(0));
  },

  toObject: function() {
    return fabric.util.object.extend(this.callSuper('toObject'), {
      cornerRadiis: this.get('cornerRadiis')
    });
  },

  _render: function(ctx) {
	var rs = this.get('cornerRadiis');
    var r = rs[0];
    for (var i=0; i<4; i++) {
      if (rs[i] != r) {
        //console.debug("drawing our own rounded rect");
        if (this.width === 1 && this.height === 1) {
          ctx.fillRect(0, 0, 1, 1);
          return;
        }

        var w = this.width,
            h = this.height,
            x = -w / 2,
            y = -h / 2,
            isInPathGroup = this.group && this.group.type === 'path-group',
            isRounded = true;

        ctx.beginPath();
        ctx.globalAlpha = isInPathGroup ? (ctx.globalAlpha * this.opacity) : this.opacity;

        if (this.transformMatrix && isInPathGroup) {
          ctx.translate(
            this.width / 2 + this.x,
            this.height / 2 + this.y);
        }
        if (!this.transformMatrix && isInPathGroup) {
          ctx.translate(
            -this.group.width / 2 + this.width / 2 + this.x,
            -this.group.height / 2 + this.height / 2 + this.y);
        }

        gosuiCanvasRRect(ctx, x, y, w, h, this.cornerRadiis);

        this._renderFill(ctx);
        this._renderStroke(ctx);
		return;        
      }
    }
	//console.debug("drawing fabric rounded rect");
    //Call the fabricjs render method if all 4 corners are the same
    this.set('rx', r);
    this.set('ry', r);
    this.callSuper('_render', ctx);
  }
});

function fabricCanvasResize(w, h) {
	canvas.setWidth(w)
	canvas.setHeight(h)
}

function fabricDrawRect(spec) {
	var rect = new RoundedRect(spec);
	//console.debug(spec);
  	gosuiAddUnsel(canvas, rect);
}