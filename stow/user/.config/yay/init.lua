yay.create_autocmd("UpgradeSelect", {
  desc = "skip AUR upgrades modified in the last day",

  callback = function(event)
    local exclude = {}
    local cutoff = os.time() - (24 * 60 * 60)

    for _, pkg in ipairs(event.data.upgrades) do
      if pkg.repository == "aur" and pkg.last_modified >= cutoff then
        local available_at = pkg.last_modified + (24 * 60 * 60)

        yay.log.warn(
          "pre-excluding recently modified AUR package: ",
          pkg.name,
          " (available " .. os.date("%Y-%m-%d %H:%M", available_at) .. ")"
        )

        table.insert(exclude, pkg.name)
      end
    end

    return {
      exclude = exclude,
      skip_menu = false,
    }
  end,
})
