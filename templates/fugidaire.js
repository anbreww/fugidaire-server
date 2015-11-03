updater = _.throttle(function(id, color, override) {
  color = "#"+color;
  console.log(id + ': '+color);
  $.ajax({
    url:'/api/v1/'+id,
    type:'PUT',
    data:{'color':color, 'apikey':getUrlParameter('apikey')},
  });
}, 500);

fan_update = function(id, new_value) {
  $.ajax({
    url:'/api/v1/fan/'+id,
    type:'PUT',
    data:{'on':new_value, 'apikey':getUrlParameter('apikey')},
  });
};

$('.checkbox-fan').on('change', function() {
  console.log(this.checked, this.id);
  fan_update(this.id, this.checked);
});

var getUrlParameter = function getUrlParameter(sParam) {
    var sPageURL = decodeURIComponent(window.location.search.substring(1)),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : sParameterName[1];
        }
    }
};

$('#apikey').val(getUrlParameter('apikey'));

$('.color').colorPicker({
  buildCallback: function($elm) {
      $elm.prepend('<div class="cp-color-picker"></div>');
  },
  cssAddon: '.cp-disp{padding:10px; margin-bottom:6px; font-size:19px; height:20px; line-height:20px}' +
  '.cp-xy-slider{width:200px; height:200px;}' +
'.cp-xy-cursor{width:16px; height:16px; border-width:2px; margin:-8px}' +
            '.cp-z-slider{height:200px; width:40px;}' +
            '.cp-z-cursor{border-width:8px; margin-top:-8px;}' +
            '.cp-alpha{height:40px;}' +
            '.cp-alpha-cursor{border-width: 8px; margin-left:-8px;}',
  opacity: false,
  doRender: 'div div',
  preventFocus: true,
  renderCallback: function($elm, toggled) {
    var colors = this.color.colors,
    rgb = colors.RND.rgb;
    updater($elm.attr('id'), colors.HEX, false);

    if (toggled === true) {
      // here you can recalculate position after showing the color picker
      // in case it doesn't fit into view.
      $('.trigger').removeClass('active');
      $elm.closest('.trigger').addClass('active');
    } else if (toggled === false) {
      // this happens when the box just disappeared
      $elm.closest('.trigger').removeClass('active');
    }

    $('.cp-disp').css({
        backgroundColor: '#' + colors.HEX,
        color: colors.RGBLuminance > 0.22 ? '#222' : '#ddd'
    }).text('rgba(' + rgb.r + ', ' + rgb.g + ', ' + rgb.b +
        ', ' + (Math.round(colors.alpha * 100) / 100) + ')');
    },
}); // that's it
