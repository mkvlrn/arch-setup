// ==UserScript==
// @name        enable proton login auto complete
// @namespace   Violentmonkey Scripts
// @icon        https://meurh.spread.com.br/teleinformatica/assets/img/favicon-32x32.png
// @version     1.0.0
//
// @match       https://meurh.spread.com.br/teleinformatica/*
// @grant       none
//
// @author      dev@mkvlrn.cc
// @description
// ==/UserScript==

new MutationObserver(() => {
  const user = document.querySelector('input[name="user"]');
  const pass = document.querySelector('input[name="password"]');

  if (user) user.setAttribute("autocomplete", "username");
  if (pass) pass.setAttribute("autocomplete", "current-password");
}).observe(document.documentElement, {
  childList: true,
  subtree: true,
});
