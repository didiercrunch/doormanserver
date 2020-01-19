(function(){
  const Ratio = (n, d) => ({
          type: 'Ratio',
          'n': n, // numerator
          'd': d // denominator
      });

  const approxRatio = eps => n => {
        const
            gcde = (e, x, y) => {
                const _gcd = (a, b) => (b < e ? a : _gcd(b, a % b));
                return _gcd(Math.abs(x), Math.abs(y));
            },
            c = gcde(Boolean(eps) ? eps : (1 / 10000), 1, n);
        return Ratio(
            Math.floor(n / c), // numerator
            Math.floor(1 / c) // denominator
        );
    };

  var topBar = {}

  function fontawesome(){
    const classes = _.map(arguments, x => "." + x).reduce((s, n) => s + n);
    return m('i' + classes);
  }

  var TopBar = {
    view: function(){
      return m('nav.navbar', [
        m('div.navbar-brand',
          m('a.navbar-item', {href: "#!/"}, "Doorman ", fontawesome('fas', 'fa-door-open'))),
        m('div.navbar-menu', [
          m('div.navbar-end', [
            m('a.navbar-item',  {href: "#!/server"}, "Server"),
            m('a.navbar-item',  {href: "#!/new-doorman"}, "New Doorman"),
          ])
        ])
      ]);
    }
  }

  var Home = {};

  Home.doormen = [];

  Home.oninit = function(){
    m.request({method: "GET", url: "/api/doormen"}).then(data => Home.doormen = data.doormen);
  }

  Home.doorman = function(doorman, idx){
    const link = "#!/doorman/" + doorman.id
    return m('li', m('a', {href: link}, doorman.name));
  }

  Home.view = function() {
      return m("main", [
          m(TopBar),
          m("h1", {class: "title"}, "Home"),
          m('div.content',
            m('ul', _.map(Home.doormen, Home.doorman)))
      ]);
  }

  var NewDoorman = {}


  NewDoorman.oninit = function(){
    NewDoorman.clear();
  }

  NewDoorman.setToPayload = function(field){
    return function(e){
      _.set(NewDoorman.payload, field, e.target.value);
    }
  }

  NewDoorman.nameField = function(){
    const fields = {
      type: 'text',
      placeholder: 'Background Colour dsds',
      onchange: NewDoorman.setToPayload('foo')};
    return m('label.label', 'Name',
       m('input.input', fields));
  }

  NewDoorman.ownerField = function(email, idx){
    const label = 'Owner ' + (idx + 1);
    const fields = {
      type: 'text',
      value: email,
      onchange: NewDoorman.setToPayload(['OwnerEmails', idx])}
    return m('label.label', label,
      m('input.input', fields));
  }

  NewDoorman.putEmptyOwnerFields = function(){
    NewDoorman.payload.OwnerEmails = _(NewDoorman.payload.OwnerEmails)
      .filter(x => !_.isEmpty(x))
      .concat('')
      .value();
  }

  NewDoorman.ownerFields = function(){
    return m('fieldset', {onchange: NewDoorman.putEmptyOwnerFields()},
      m('legend', 'Owners'),
        _.map(NewDoorman.payload.OwnerEmails, NewDoorman.ownerField));
  }
  NewDoorman.valueField = function(value, idx){
    const nameField = m('div.field',
      m('label', 'name',
        m('input.input', {type: 'text', value: value.name, onchange: NewDoorman.setToPayload(['Values', idx, 'name'])})));
    const valueField = m('div.field', m('label', 'value', m('input.input', {type: 'number', value: value.value, onchange: NewDoorman.setToPayload(['Values', idx, 'value'])})));

    return m('div.columns',
        m('div.column', nameField),
        m('div.column', valueField));
  }

  NewDoorman.putEmptyValueFields = function() {
    const hasEmptyNames = _(NewDoorman.payload.Values).map('name').filter(_.isEmpty).value().length;
    const hasEmptyValues = _(NewDoorman.payload.Values).map('value').filter(_.isEmpty).value().length;
    const sum = _(NewDoorman.payload.Values).map('value').map(x => parseInt(x)).sum();
    if(sum >= 100 || hasEmptyNames || hasEmptyValues){
      return;
    }
    NewDoorman.payload.Values.push({name: '', value: 100 - sum});
  }

  NewDoorman.valueFields = function(){
    return m('fieldset', {onchange: NewDoorman.putEmptyValueFields}, m('legend', 'Values'),
      _.map(NewDoorman.payload.Values, NewDoorman.valueField));
  }

  NewDoorman.clear = function(){
    NewDoorman.payload = {
      OwnerEmails: [],
      Values: [{}]
    };
  }

  NewDoorman.buttons = function(){
    return m('div.buttons', [
      m('button.button.is-primary.is-light', 'Create'),
      m('button.button.is-danger.is-light', {onclick: NewDoorman.clear}, 'Clear')
    ])
  }

  NewDoorman.view = function(){
    return m('div',
      m(TopBar),
      m('form', [
        NewDoorman.nameField(),
        m('br'),
        NewDoorman.ownerFields(),
        m('br'),
        NewDoorman.valueFields(),
        m('br'),
        NewDoorman.buttons()
      ]))
  }

  var Server = {}

  Server.data = {};

  Server.oninit= function(){
    m.request({method: "GET", url: "/api/server"})
    .then(function(data) {
        Server.data = data;
    });
  }

  Server.renderRow = function(value, key){
    return m('tr', [
      m('td', key),
      m('td', _.isObject(value) ? JSON.stringify(value) : value)]);
  }

  Server.renderTable = function(serverData){
    const rows = _.map(serverData, Server.renderRow);
    return m('table.table.is-bordered.is-striped.is-hoverable', rows);
  }

  Server.view = function() {
    return m("div", [
             m(TopBar),
             m("h1", {class: "title"}, "Server"),
             Server.renderTable(Server.data)
      ]);
  }

  var Doorman = {};
  Doorman.doorman = {}

  Doorman.oninit = function(vnode){
    m.request({method: "GET", url: "/api/doormen"})
      .then(data => _(data.doormen).filter(x => x.id == vnode.attrs.doormanId).map('url').first())
      .then(url => m.request({method: "GET", url: url}))
      .then(data => Doorman.doorman = data);
  }

  Doorman.view = function(vnode){
    return m('div', [m(TopBar),'blah' + JSON.stringify(Doorman.doorman)])

  }

  m.route(document.getElementById('mithril'), "/", {
    "/": Home,
    "/server": Server,
    "/new-doorman": NewDoorman,
    "/doorman/:doormanId": Doorman,
  });

})();
